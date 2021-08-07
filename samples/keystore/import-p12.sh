#!/bin/bash -x

keypath=$1
certpath=$2
outdir=$3

p12path="$outdir/client_key.p12"
kdbpath="$outdir/key.kdb"

certlabel="ibmwebspheremqkarson"
#certlabel="cert"

/usr/bin/openssl pkcs12 -export -name $certlabel -out $p12path -inkey $keypath -in $certpath -keypbe NONE -certpbe NONE -nomaciter -passout "pass:"

/opt/mqm/bin/runmqckm -cert -import -file $p12path -pw "" -type pkcs12 -target $kdbpath -target_stashed -target_type cms -label $certlabel -new_label $certlabel

rm $p12path
