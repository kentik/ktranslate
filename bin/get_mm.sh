#!/bin/bash

curl -o mm.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=${MM_DOWNLOAD_KEY}&suffix=tar.gz"
tar -zxf mm.tar.gz
mv GeoLite2-Country_*/GeoLite2-Country.mmdb config/GeoLite2-Country.mmdb
rm mm.tar.gz
rm -r GeoLite2-Country_*
