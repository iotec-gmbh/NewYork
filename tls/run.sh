#!/bin/bash

set -o errexit
set -o nounset


echo "CERTBOT_WEBROOT=$CERTBOT_WEBROOT
CERTBOT_EMAIL=$CERTBOT_EMAIL
CERTBOT_DOMAINS=$CERTBOT_DOMAINS
CERTBOT_TEST_CERT=${CERTBOT_TESTCERT}
"

domain_args=""
for domain in $CERTBOT_DOMAINS
do
    domain_args="$domain_args --domain $domain"
done

mkdir -p "${CERTBOT_WEBROOT}/.well-known/acme-challenge"

test_cert=""
if [ "$CERTBOT_TESTCERT" = "y" ]
then
    test_cert="--test-cert"
fi

for domain in $CERTBOT_DOMAINS
do
    certbot certonly $test_cert --webroot -w "$CERTBOT_WEBROOT" --agree-tos -m "$CERTBOT_EMAIL" --non-interactive --domain "$domain"
done

while true
do
    sleep 1d
    certbot --test-cert renew
done
