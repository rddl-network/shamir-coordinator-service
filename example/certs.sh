#!/bin/bash
set -o errexit -o errtrace

CWD=${PWD##*/}
if [[ "$CWD" != "example" ]]; then
    echo "err: please run script from example directory!"
    exit 1
fi

# output directory
mkdir -p certs
pushd certs >/dev/null

CANAME=ca
FILE="$CANAME.key"
if [[ -f "$FILE" ]]; then
    echo "err: $FILE already exists!"
    exit 2
fi

# generate ed25519 private key
openssl genpkey -algorithm ed25519 -out $CANAME.key
# create certificate, 3650 days = 10 years
openssl req -x509 -new -nodes -key $CANAME.key -sha256 -days 3650 -out $CANAME.crt -subj '/CN=localhost'

for MYCERT in server client; do
    # create certificate signing request
    openssl req -new -nodes -out $MYCERT.csr -newkey ed25519 -keyout $MYCERT.key -subj '/CN=localhost'
    # create certificate, 730 days = 2 years
    openssl x509 -req -in $MYCERT.csr -CA $CANAME.crt -CAkey $CANAME.key -out $MYCERT.crt -days 730 -sha256 -extfile ../$MYCERT.v3.ext
    # remove certificate signing request
    rm $MYCERT.csr
done

# sonity check
ls -l
popd >/dev/null
