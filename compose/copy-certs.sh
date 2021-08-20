#!/bin/bash -x

tlsdir=$1

if [[ -z $tlsdir ]]; then
echo tls directory required: copy-certs.sh '<tls-gen-dir>'. Use tls-gen to create keys and certs and point to the result directory.
exit 1
fi

cp $tlsdir/server_key_nopass.pem output/etc/mqm/pki/cert/tls.key
cp $tlsdir/server_certificate.pem output/etc/mqm/pki/cert/tls.crt
cp $tlsdir/ca_certificate.pem output/etc/mqm/pki/cert/ca.crt

cp $tlsdir/ca_certificate.pem output/etc/mqm/pki/trust/ca_certificate.crt
