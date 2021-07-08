#!/bin/bash

# the directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# the temp directory used, within $DIR
# omit the -p parameter to create a temporal directory in the default location
WORK_DIR=`mktemp -d -p "$DIR"`

# check if tmp dir was created
if [[ ! "$WORK_DIR" || ! -d "$WORK_DIR" ]]; then
  echo "Could not create temp dir"
  exit 1
fi

# deletes the temp directory
function cleanup {
  rm -rf "$WORK_DIR"
  echo "Deleted temp working directory $WORK_DIR"
}

# register the cleanup function to be called on the EXIT signal
trap cleanup EXIT

# start testing
pushd $WORK_DIR
docker pull kentik/ktranslate:staging
cd .
id=$(docker create kentik/ktranslate:staging)
docker cp $id:/etc/ktranslate/snmp-base.yaml .
docker rm -v $id
docker run -ti --name ktranslate --rm --net=host \
  -v `pwd`/snmp-base.yaml:/snmp-base.yaml \
  kentik/ktranslate:staging \
    -snmp /snmp-base.yaml \
    -log_level info \
    -snmp_discovery=true
cat snmp-base.yaml
popd
