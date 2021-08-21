#!/bin/bash -x

envfile=$1

if [[ -z $envfile ]]; then
echo env file required: run-ldap-container.sh '<envfile>'
exit 1
fi

# load env
. $envfile

LDAP_ALLOW_ANON_BINDING=no

env="-e LDAP_ROOT=$LDAP_ROOT -e LDAP_ADMIN_USERNAME=$LDAP_ADMIN_USERNAME -e LDAP_ADMIN_PASSWORD=$LDAP_ADMIN_PASSWORD -e LDAP_ALLOW_ANON_BINDING=$LDAP_ALLOW_ANON_BINDING"

name=openldap
net=qmnet

sudo podman run --rm --name $name --network $net -v ldif:/ldifs $env -p $LDAP_TCP_PORT:1389 -p $LDAP_SSL_PORT:1636 docker.io/bitnami/openldap:latest
