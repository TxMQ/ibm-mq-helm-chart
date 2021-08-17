#!/bin/bash -x

qmenv=$1
qmname=$2

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

# ldap
LDAP_BIND_PASSWORD=admin
EOF
