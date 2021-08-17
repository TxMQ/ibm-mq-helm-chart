#!/bin/bash -x

tlsdir=$1

cp $tlsdir/server_key_nopass.pem output/etc/mqm/pki/cert/tls.key
cp $tlsdir/server_certificate.pem output/etc/mqm/pki/cert/tls.crt
cp $tlsdir/ca_certificate.pem output/etc/mqm/pki/cert/ca.crt

cp $tlsdir/ca_certificate.pem output/etc/mqm/pki/trust/ca_certificate.crt
