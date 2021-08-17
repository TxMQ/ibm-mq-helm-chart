#!/bin/bash -x

# read env
. ../env.sh

tag=$MQIMGTAG
img=$DC_MQIMGREG/txmq-mq-base-rpm-$MQVER:$tag

qmname=${1:-"qm1"}
net=qmnet

# vault
vault="-e VAULT_ENABLE_TLS=false -e VAULT_LDAP_CREDS_INJECT_PATH= -e VAULT_TLS_KEY_INJECT_PATH= -e VAULT_TLS_CERT_INJECT_PATH= -e VAULT_TLS_CA_INJECT_PATH="

# git
git="-e GIT_CONFIG_URL= -e GIT_CONFIG_REF= -e GIT_CONFIG_DIR="

# web
web="-e MQ_START_MQWEB=1"

# debug
debug="-e MQRUNNER_DEBUG=1"

# tls
tls="-e MQ_ENABLE_TLS_NO_VAULT=1"

# ldap
ldap="-e LDAP_BIND_PASSWORD=admin"

# qmgr, required
qmgr="-e MQ_QMGR_NAME=$qmname"

# all envs
envars="$qmgr $debug $tls $web $vault $git $ldap"

# run
sudo podman run --rm --name $qmname --network $net -v mqdata:/var/mqm -v mqsc:/etc/mqm/mqsc -v qmtls:/etc/mqm/pki/cert -v qmtrust:/etc/mqm/pki/trust -v webuser:/etc/mqm/webuser $envars -p 1414:1414 -p 9443:9443 -p 40000:40000 $img
