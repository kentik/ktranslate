# build ktranslate
FROM golang:1.19-alpine as build
RUN apk add -U libpcap-dev alpine-sdk bash libcap
COPY . /src
WORKDIR /src
ARG KENTIK_KTRANSLATE_VERSION
RUN make

# maxmind dbs
FROM alpine:latest as maxmind
ARG MAXMIND_LICENSE_KEY
RUN apk add -U curl tar
ENV GEOLITE2_COUNTRY_FILE=GeoLite2-Country.mmdb
ENV GEOLITE2_ASN_FILE=GeoLite2-ASN.mmdb
RUN if [ -z "${MAXMIND_LICENSE_KEY}" ]; then echo "MAXMIND_LICENSE_KEY" not set; exit 1; fi
RUN curl -o /tmp/country.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" && \
	tar zxf /tmp/country.tar.gz --strip-components 1 -C /
RUN curl -o /tmp/asn.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz" && \
	tar zxf /tmp/asn.tar.gz --strip-components 1 -C /

# snmp profiles
FROM alpine:latest as snmp
RUN apk add -U git

# If there is a branch of snmp-profiles to use, switch over here now.
RUN if [ -z "${KENTIK_SNMP_PROFILE_REPO}" ]; then \
    git clone https://github.com/kentik/snmp-profiles /snmp; \
else \
    echo "picking repo ${KENTIK_SNMP_PROFILE_REPO} for snmp profiles"; \
    git clone ${KENTIK_SNMP_PROFILE_REPO} /snmp; \
fi

# main image
FROM alpine:3.14
RUN apk add -U --no-cache ca-certificates libpcap
RUN addgroup -g 1000 ktranslate && \
	adduser -D -u 1000 -G ktranslate -H -h /etc/ktranslate ktranslate
#RUN set -eux; \
#	groupadd --gid 1000 ktranslate; \
#	useradd --home-dir /etc/ktranslate --gid ktranslate --no-create-home --uid 1000 ktranslate

COPY --chown=ktranslate:ktranslate config/ /etc/ktranslate/
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
