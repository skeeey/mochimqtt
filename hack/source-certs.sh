# !/bin/bash

CURRENT_DIR="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"
REPO_DIR="$(cd $(dirname ${CURRENT_DIR}) && pwd)"

CERT_DIR=${REPO_DIR}/hack/certs
SOURCE_DIR=${CERT_DIR}/source

namespace="mochi-mqtt"
service="mochi-mqtt"
route="mochi-mqtt"
domain=apps.server-foundation-sno-r8b9r.dev04.red-chesterfield.com

source_name="cluster-maestro"

rm -rf $SOURCE_DIR
mkdir $SOURCE_DIR

# sign source cert and key

cat > ${SOURCE_DIR}/client-csr.json <<EOF
{
  "CN":"${source_name}",
  "key":{
    "algo":"rsa",
    "size":2048
  },
  "hosts":[
    "127.0.0.1",
    "localhost",
    "${service}",
    "${service}.${namespace}",
    "${service}.${namespace}.svc",
    "${service}.${namespace}.svc.cluster.local",
    "${service}-${route}.${domain}"
  ],
  "names":[
    {
      "C":"CN",
      "ST":"ShaaXi",
      "L":"Xi'an",
      "O":"ACM",
      "OU":"core"
    }
  ]
}
EOF

cfssl gencert -ca=${CERT_DIR}/root-ca.pem -ca-key=${CERT_DIR}/root-ca-key.pem \
  -config=${CERT_DIR}/profiles.json \
  -profile=test ${SOURCE_DIR}/client-csr.json | cfssljson -bare ${SOURCE_DIR}/client
