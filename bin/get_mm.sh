#!/usr/bin/env bash

if [ -z ${MAXMIND_LICENSE_KEY} ]; then
  echo ""
  echo "Provide a MaxMind license key to fetch the Geolite2 database files."
  echo "export MAXMIND_LICENSE_KEY=<your-license-key>"
  echo ""
  echo "Sign-up for a free GeoLite2 account at https://dev.maxmind.com and create your license key."
  echo ""
  exit -1
fi

mkdir -p config
curl -L -o mm.tar.gz -u ${YOUR_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY} "https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz"
tar -zxf mm.tar.gz
mv GeoLite2-Country_*/GeoLite2-Country.mmdb config/GeoLite2-Country.mmdb
rm mm.tar.gz
rm -r GeoLite2-Country_*

curl -L -o mm.tar.gz -u ${YOUR_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY} "https://download.maxmind.com/geoip/databases/GeoLite2-ASN/download?suffix=tar.gz"
tar -zxf mm.tar.gz
mv GeoLite2-ASN_*/GeoLite2-ASN.mmdb config/GeoLite2-ASN.mmdb
rm mm.tar.gz
rm -r GeoLite2-ASN_*
