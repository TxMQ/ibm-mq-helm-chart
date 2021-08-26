#!/bin/bash -x

. ../env.sh

mkdir -p output

qmname=${1:-qm1}

ldaptype=${LDAP_TYPE:openldap}
ldapns=${LDAP_NS:-default}
ldaphost=${LDAP_HOST:-openldap.$ldapns.svc.cluster.local}
ldapport=${LDAP_PORT:-389}
ldaproot=${LDAP_ROOT:-dc=mqldap,dc=com}
ldapuser=${LDAP_USER:-cn=admin,$ldaproot}
basednu=${BASEDN_USERS:-ou=users,dc=$ldaproot}
basedng=${BASEDN_GROUPS:-ou=groups,$ldaproot}

cat <<EOF > output/$qmname.env
# queue manager name, do not change
QMNAME=$qmname

# mq immage params, do not change
MQIMGTAG=$MQIMGTAG
MQVER=$MQVER

# mq image registry
MQIMGREG=$MQIMGREG

# ldap params
LDAP_TYPE=$ldaptype
LDAP_HOST=$ldaphost
LDAP_PORT=$ldapport
LDAP_ROOT=$ldaproot
LDAP_USER=$ldapuser
BASEDN_USERS=$basednu
BASEDN_GROUPS=$basedng
EOF
