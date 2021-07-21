FROM ubuntu:20.04

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libpcap0.8 \
	; \
	rm -rf /var/lib/apt/lists/*

COPY config/ /etc/ktranslate/

# add backwards compatibility symlinks for folks using an snmp.yml from the older image (and "ls" to verify the symlinks are correct and working)
RUN set -eux; ln -sv ktranslate/profiles ktranslate/mibs.db /etc/; ls -l /etc/profiles/ /etc/mibs.db/

COPY bin/ktranslate /usr/local/bin/

ENTRYPOINT ["ktranslate", "-metalisten", "0.0.0.0:8083", "-listen", "0.0.0.0:8082", "-mapping", "/etc/ktranslate/config.json", "-geo", "/etc/ktranslate/GeoLite2-Country.mmdb", "-udrs", "/etc/ktranslate/udr.csv", "-api_devices", "/etc/ktranslate/devices.json", "-asn", "/etc/ktranslate/GeoLite2-ASN.mmdb", "-log_level", "info"]

EXPOSE 8082
EXPOSE 8083
