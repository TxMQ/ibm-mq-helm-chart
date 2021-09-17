#!/bin/bash -x

mkdir -p output

qmname=${1:-"qm1"}

# network
./create-network.sh

# ldap container environment
ldapenv=output/ldap.env
if [[ ! -f $ldapenv ]]; then
./set-ldap-container-env.sh $ldapenv
fi

# mq container environment
qmenv=output/$qmname.env
if [[ ! -f $qmenv ]]; then
./set-mq-container-env.sh $qmenv $qmname
fi

# mqldap
if [[ ! -d output/etc/mqm/mqyaml ]]; then
mkdir -p output/etc/mqm/mqyaml
./mqldap-template.sh $qmenv
fi

# webuser
if [[ ! -d output/etc/mqm/webuser ]]; then
mkdir -p output/etc/mqm/webuser
./webuser-template.sh $qmenv
fi

# mqsc
if [[ ! -d output/etc/mqm/mqsc ]]; then
mkdir -p output/etc/mqm/mqsc
./mqscic-template.sh $qmenv
./mqexplorer-mqsc-template.sh $qmenv
fi

# qmini
if [[ ! -d output/etc/mqm/qmini ]]; then
mkdir -p output/etc/mqm/qmini
./qmini-template.sh $qmenv
fi

# ldif
if [[ ! -d output/ldif ]]; then
mkdir output/ldif
./ldif-template.sh $ldapenv
fi

# tls certs
if [[ ! -d output/etc/mqm/pki ]]; then
mkdir -p output/etc/mqm/pki/cert
mkdir -p output/etc/mqm/pki/trust
if [[ ! -z $TLS_GEN_RESULT ]]; then
./copy-certs.sh $TLS_GEN_RESULT
else
echo TLS_GEN_RESULT env var not set, certificates not copied\; copy certificates manually: copy-certs.sh '<tls-gen-result-dir>'
fi
fi

# docker compose
if [[ ! -f output/docker-compose.yaml ]]; then
./docker-compose-template.sh $ldapenv $qmenv
fi
