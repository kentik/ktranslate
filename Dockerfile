FROM ubuntu:20.04
RUN apt-get update && apt-get install -y \
    krb5-user libsasl2-modules-gssapi-mit liblz4-dev libzstd-dev libsasl2-dev libpcap-dev ca-certificates wget zstd
COPY bin/ktranslate /usr/bin/ktranslate
COPY bin/fetch /usr/bin/fetch
COPY config.json /etc/config.json
COPY code2city.mdb /etc/code2city.mdb
COPY code2region.mdb /etc/code2region.mdb
COPY ipv4-asn-ip.csv /etc/ipv4-asn-ip.csv
COPY ipv6-asn-ip.csv /etc/ipv6-asn-ip.csv
COPY asn-to-name.tsv /etc/asn-to-name.tsv
COPY udr.csv /etc/udr.csv
COPY lib/librdkafka.so.1 /lib/x86_64-linux-gnu/librdkafka.so.1
COPY mibs.db /etc/mibs.db
COPY devices.json /etc/devices.json
COPY profiles /etc/profiles
COPY snmp.yaml.sample /etc/snmp.yaml.sample
COPY snmp-base.yaml /etc/snmp-base.yaml
COPY devices.json /etc/devices.json
ENTRYPOINT ["/usr/bin/ktranslate", "-metalisten", "0.0.0.0:8083", "-metrics", "tsdb:https://flow.kentik.com/tsdb", "-listen", "0.0.0.0:8082", "-mapping", "/etc/config.json", "-city", "/etc/code2city.mdb", "-region", "/etc/code2region.mdb", "-udrs", "/etc/udr.csv", "-api_devices", "/etc/devices.json", "-asn4", "/etc/ipv4-asn-ip.csv", "-asn6", "/etc/ipv6-asn-ip.csv", "-asnName", "/etc/asn-to-name.tsv", "-log_level", "info"]
WORKDIR /

EXPOSE 8082
EXPOSE 8083
