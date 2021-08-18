#!/bin/bash
# Called by go generate

ver="${KENTIK_KTRANSLATE_VERSION:-}"
if [[ "$ver" = "" ]]; then
	ver=`git rev-parse HEAD`
	if test -n "$(git status --porcelain)"; then
		ver="dirty-$ver"
	fi
	date="un-dated build"
else
	date=`date`
fi
platform=`uname -srm`
golang=`go version`

tmpfile=`mktemp`
sed -n "/BEGIN-TEMPLATE/,/END-TEMPLATE/p" ${GOFILE} | grep -v BEGIN-TEMPLATE | grep -v END-TEMPLATE \
	| sed -e "s/\@VERSION@/${ver}/g" \
	| sed -e "s/\@DATE@/${date}/g" \
	| sed -e "s|\@PLATFORM@|${platform}|g" \
	| sed -e "s|\@GOLANG@|${golang}|g" \
>${tmpfile}

if ! cmp -s ${tmpfile} generated_${GOFILE}; then
	mv ${tmpfile} generated_${GOFILE}
fi
