# !/bin/bash

CURRENT_DIR="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"
REPO_DIR="$(cd $(dirname ${CURRENT_DIR}) && pwd)"

certpath="$REPO_DIR/hack/certs"

$REPO_DIR/bin/server -ca-path=$certpath/root-ca.pem \
    -server-cert-path=$certpath/server.pem\
    -server-key-path=$certpath/server-key.pem\
