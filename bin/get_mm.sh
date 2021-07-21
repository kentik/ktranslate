#!/bin/bash

curl -o mm.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=${MM_DOWNLOAD_KEY}&suffix=tar.gz"
tar -zxf mm.tar.gz
mv GeoLite2-Country_*/GeoLite2-Country.mmdb config/GeoLite2-Country.mmdb
rm mm.tar.gz
rm -r GeoLite2-Country_*

curl -o mm.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=${MM_DOWNLOAD_KEY}&suffix=tar.gz"
tar -zxf mm.tar.gz
mv GeoLite2-ASN_*/GeoLite2-ASN.mmdb config/GeoLite2-ASN.mmdb
rm mm.tar.gz
rm -r GeoLite2-ASN_*
