#!/bin/bash -x

. ../env.sh

mkdir -p output

envfile=$1
if [[ -z $envfile ]]; then
echo env file name required: set-ldap-container-env.sh '<envfile>'
exit 1
fi

ldaproot=${LDAP_ROOT:-dc=mqldap,dc=com}
adminusername=${LDAP_ADMIN_USERNAME:-admin}
adminpassword=${LDAP_ADMIN_PASSWORD:-admin}
tcpport=${LDAP_TCP_PORT:-1389}
sslport=${LDAP_SSL_PORT:-1636}

# write out env file
cat <<EOF > $envfile
# ldap container params
LDAP_ROOT=$ldaproot
LDAP_ADMIN_USERNAME=$adminusername
LDAP_ADMIN_PASSWORD=$adminpassword
LDAP_TCP_PORT=$tcpport
LDAP_SSL_PORT=$sslport
EOF
