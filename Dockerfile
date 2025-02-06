# build ktranslate
FROM golang:1.23-alpine as build
RUN apk add -U libpcap-dev alpine-sdk bash libcap
COPY . /src
WORKDIR /src
ARG KENTIK_KTRANSLATE_VERSION
RUN make

# maxmind dbs
FROM alpine:latest as maxmind
ARG MAXMIND_LICENSE_KEY
ARG YOUR_ACCOUNT_ID
RUN apk add -U curl tar
ENV GEOLITE2_COUNTRY_FILE=GeoLite2-Country.mmdb
ENV GEOLITE2_ASN_FILE=GeoLite2-ASN.mmdb
RUN if [ -z "${MAXMIND_LICENSE_KEY}" ]; then echo "MAXMIND_LICENSE_KEY" not set; exit 1; fi
RUN curl -L -o /tmp/country.tar.gz -u ${YOUR_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY} "https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz" && \
	tar zxf /tmp/country.tar.gz --strip-components 1 -C /
RUN curl -L -o /tmp/asn.tar.gz -u ${YOUR_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY} "https://download.maxmind.com/geoip/databases/GeoLite2-ASN/download?suffix=tar.gz" && \
	tar zxf /tmp/asn.tar.gz --strip-components 1 -C /

# snmp profiles
FROM alpine:latest as snmp
ARG KENTIK_SNMP_PROFILE_REPO
RUN apk add -U git

# If there is a branch of snmp-profiles to use, switch over here now.
RUN if [ -z "${KENTIK_SNMP_PROFILE_REPO}" ]; then \
    git clone https://github.com/kentik/snmp-profiles /snmp; \
else \
    echo "picking repo ${KENTIK_SNMP_PROFILE_REPO} for snmp profiles"; \
    git clone ${KENTIK_SNMP_PROFILE_REPO} /snmp; \
fi

# main image
FROM alpine:3.21
RUN apk add -U --no-cache ca-certificates libpcap aws-cli
RUN addgroup -g 1000 ktranslate && \
	adduser -D -u 1000 -G ktranslate -H -h /etc/ktranslate ktranslate
#RUN set -eux; \
#	groupadd --gid 1000 ktranslate; \
#	useradd --home-dir /etc/ktranslate --gid ktranslate --no-create-home --uid 1000 ktranslate

# Some people want to specify an alternative config dir. This lets them override with --build-arg CONFIG-DIR=my-new-dir
ARG CONFIG_DIR=config
COPY --chown=ktranslate:ktranslate ${CONFIG_DIR}/ /etc/ktranslate/
COPY --chown=ktranslate:ktranslate lib/ /etc/ktranslate/

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
ENTRYPOINT ["ktranslate", "-listen", "off", "-mapping", "/etc/ktranslate/config.json", "-geo", "/etc/ktranslate/GeoLite2-Country.mmdb", "-udrs", "/etc/ktranslate/udr.csv", "-api_devices", "/etc/ktranslate/devices.json", "-asn", "/etc/ktranslate/GeoLite2-ASN.mmdb", "-log_level", "info"]
