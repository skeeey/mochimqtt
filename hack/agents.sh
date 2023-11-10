# !/bin/bash

CURRENT_DIR="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"
REPO_DIR="$(cd $(dirname ${CURRENT_DIR}) && pwd)"

logdir="$REPO_DIR/logs"
certpath="$REPO_DIR/hack/certs"
host="mochi-mqtt-mochi-mqtt.apps.server-foundation-sno-r8b9r.dev04.red-chesterfield.com:443"

rm -rf $logdir
mkdir -p $logdir

binpath=$REPO_DIR/bin

$binpath/client -ca-path=$certpath/root-ca.pem \
    -client-cert-path=$certpath/cluster1/client.pem\
    -client-key-path=$certpath/cluster1/client-key.pem\
    -host=$host >$logdir/$i.log 2>&1 &
