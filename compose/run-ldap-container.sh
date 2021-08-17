#!/bin/bash -x

LDAP_ROOT=dc=mqldap,dc=com
LDAP_ADMIN_USERNAME=admin
LDAP_ADMIN_PASSWORD=admin
LDAP_ALLOW_ANON_BINDING=no

env="-e LDAP_ROOT=$LDAP_ROOT -e LDAP_ADMIN_USERNAME=$LDAP_ADMIN_USERNAME -e LDAP_ADMIN_PASSWORD=$LDAP_ADMIN_PASSWORD -e LDAP_ALLOW_ANON_BINDING=$LDAP_ALLOW_ANON_BINDING"

name=openldap
net=qmnet

sudo podman run --rm --name $name --network $net -v ldif:/ldifs $env -p 1389:1389 -p 1636:1636 docker.io/bitnami/openldap:latest
