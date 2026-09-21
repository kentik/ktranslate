# build ktranslate
FROM golang:1.25-alpine AS build
RUN apk add -U make bash libcap
ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
ARG NETWORK_AGENT_VERSION
ARG NETWORK_AGENT_BUILD
RUN make

# maxmind dbs
#
# Downloaded and cached by the calling workflow (see ci-build.yml / publish-release.yml),
# via actions/cache -- not inside this build at all, since a BuildKit `--mount=type=cache`
# doesn't survive across CI runs (each job gets a fresh BuildKit daemon) and isn't included
# in `--cache-from`/`--cache-to` exports, so it can never actually persist here. This stage
# just stages the already-downloaded files from the build context for the COPY below.
# Building locally: run `just maxmind-dbs` first to populate maxmind-dbs/.
FROM scratch AS maxmind
COPY maxmind-dbs/GeoLite2-Country.mmdb /GeoLite2-Country.mmdb
COPY maxmind-dbs/GeoLite2-ASN.mmdb /GeoLite2-ASN.mmdb

# snmp profiles
FROM alpine:latest AS snmp
ARG NR_SNMP_PROFILE_REPO
RUN apk add -U git

# Opt-in auth: when a `github_token` BuildKit secret is provided (a GitHub token with read
# access to the repo), transparently authenticate GitHub HTTPS clones. This is a complete
# no-op when the secret is absent, so the override/clone logic below is unchanged from
# upstream. The token lives only in this throwaway stage (only /snmp/profiles is copied on).
RUN --mount=type=secret,id=github_token \
    if [ -s /run/secrets/github_token ]; then \
        git config --global url."https://x-access-token:$(cat /run/secrets/github_token)@github.com/".insteadOf "https://github.com/"; \
    fi

# If there is a branch of snmp-profiles to use, switch over here now.
RUN if [ -z "${NR_SNMP_PROFILE_REPO}" ]; then \
    git clone https://github.com/kentik/snmp-profiles /snmp; \
else \
    echo "picking repo ${NR_SNMP_PROFILE_REPO} for snmp profiles"; \
    git clone ${NR_SNMP_PROFILE_REPO} /snmp; \
fi

# main image
FROM alpine:3.23.3
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1000 ktranslate && \
	adduser -D -u 1000 -G ktranslate -H -h /etc/ktranslate ktranslate
#RUN set -eux; \
#	groupadd --gid 1000 ktranslate; \
#	useradd --home-dir /etc/ktranslate --gid ktranslate --no-create-home --uid 1000 ktranslate

# Some people want to specify an alternative config dir. This lets them override with --build-arg CONFIG-DIR=my-new-dir
ARG CONFIG_DIR=config
COPY --chown=ktranslate:ktranslate ${CONFIG_DIR}/ /etc/ktranslate/

# maxmind db
COPY --from=maxmind /GeoLite2-Country.mmdb /etc/ktranslate/
COPY --from=maxmind /GeoLite2-ASN.mmdb /etc/ktranslate/
# snmp
COPY --from=snmp /snmp/profiles /etc/ktranslate/profiles

# add backwards compatibility symlinks for folks using an snmp.yml from the older image (and "ls" to verify the symlinks are correct and working)
RUN ls -lah /etc/ktranslate ; ln -sv /etc/ktranslate /etc/profiles ; ls -lah /etc/profiles/
RUN ln -sv /etc/ktranslate/mibs.db /etc/mib.db ; ls -lah /etc/mib.db/

COPY --from=build /src/bin/ktranslate /usr/local/bin/ktranslate
COPY --from=build /usr/sbin/setcap /usr/sbin/setcap
COPY --from=build /usr/lib/libcap.so.2 /usr/lib/libcap.so.2
RUN setcap cap_net_raw=+ep /usr/local/bin/ktranslate

EXPOSE 8082

USER ktranslate
ENTRYPOINT ["ktranslate", "-listen", "off", "-mapping", "/etc/ktranslate/config.json", "-geo", "/etc/ktranslate/GeoLite2-Country.mmdb", "-udrs", "/etc/ktranslate/udr.csv", "-api_devices", "/etc/ktranslate/devices.json", "-asn", "/etc/ktranslate/GeoLite2-ASN.mmdb", "-log_level", "info", "-geo_region_map", "/etc/ktranslate/ch_region_mapping.csv.gz", "-geo_city_map", "/etc/ktranslate/ch_city_mapping.csv.gz"]
