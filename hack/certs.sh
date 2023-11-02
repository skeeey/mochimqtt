# !/bin/bash

rm -f *.json
rm -f *.pem
rm -f *.csr

cat > profiles.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "test": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "8760h"
      }
    }
  }
}
EOF

# generate root ca and key
cat > root-ca.json <<EOF
{
  "CN": "Mochi MQTT Root CA",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "ShaaXi",
      "L": "Xi'an",
      "O": "ACM",
      "OU": "core"
    }
 ],
 "ca": {
    "expiry": "876000h"
 }
}
EOF

cfssl gencert -initca root-ca.json | cfssljson -bare root-ca


# sign server cert and key

cat > server-csr.json <<EOF
{
  "CN":"skeeey.macbook.test",
  "key":{
    "algo":"rsa",
    "size":2048
  },
  "hosts":[
    "127.0.0.1",
    "localhost",
    "skeeey.macbook.test"
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

cfssl gencert -ca=root-ca.pem -ca-key=root-ca-key.pem -config=profiles.json -profile=test server-csr.json | cfssljson -bare server

# sign client cert and key

cat > client-csr.json <<EOF
{
  "CN":"cluster1",
  "key":{
    "algo":"rsa",
    "size":2048
  },
  "hosts":[
    "127.0.0.1",
    "localhost",
    "skeeey.macbook.test"
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

cfssl gencert -ca=root-ca.pem -ca-key=root-ca-key.pem -config=profiles.json -profile=test client-csr.json | cfssljson -bare client
