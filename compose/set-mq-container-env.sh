#!/bin/bash -x

qmenv=$1
qmname=$2

# ldaptype: openldap, activedirectory
ldaptype=${LDAP_TYPE:-openldap}
ldaphost=${LDAP_HOST:-openldap}
ldapport=${LDAP_PORT:-1389}
ldaproot=${LDAP_ROOT:-dc=mqldap,dc=com}
ldapuser=${LDAP_USER:-cn=admin,$ldaproot}
basednu=${BASEDN_USERS:-ou=users,$ldaproot}
basedng=${BASEDN_GROUPS:-ou=groups,$ldaproot}

cat <<EOF > $qmenv
# queue manager
MQ_QMGR_NAME=$qmname

# filter mq log to standard output
# comma separated list of prefixes
# empty value will suppress mq output to std out
# special value NO_FILTER will output every line of mq log
# special value DEFAULT_FILTER will apply AMQ filter to mq output
MQ_LOG_FILTER=

# vault
VAULT_ENABLE_TLS=false 
VAULT_LDAP_CREDS_INJECT_PATH= 
VAULT_TLS_KEY_INJECT_PATH=
VAULT_TLS_CERT_INJECT_PATH= 
VAULT_TLS_CA_INJECT_PATH=

# git
GIT_CONFIG_URL= 
GIT_CONFIG_REF=
GIT_CONFIG_DIR=

# web
MQ_START_MQWEB=1

# debug
MQRUNNER_DEBUG=1

# tls
MQ_ENABLE_TLS_NO_VAULT=1

# ldap connection for webuser
LDAP_TYPE=$ldaptype
LDAP_HOST=$ldaphost
LDAP_PORT=$ldapport
LDAP_ROOT=$ldaproot
LDAP_USER=$ldapuser

# ldap connection for authinfo
BASEDN_USERS=$basednu
BASEDN_GROUPS=$basedng

# ldap password
LDAP_BIND_PASSWORD=${LDAP_BIND_PASSWORD:-admin}

# application group names
APPL_GROUP="devs"
ADMIN_GROUP="admins"
READ_ADMIN_GROUP="readadmins"
EOF
