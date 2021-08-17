#!/bin/bash -x

mkdir -p output

qmname=${1:-"qm1"}

./docker-compose-template.sh output $qmname

qmenv=output/$qmname.env
if [[ ! -f $qmenv ]];then
./qm-env-template.sh $qmenv $qmname
fi

if [[ ! -d output/etc ]]; then
cp -r etc output
mkdir -p output/etc/mqm/pki/cert
mkdir -p output/etc/mqm/pki/trust
fi

if [[ ! -d output/ldif ]]; then
cp -r ldif output
fi
