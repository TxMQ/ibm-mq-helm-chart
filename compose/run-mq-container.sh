#!/bin/bash -x

# read env
. ../env.sh

tag=$MQIMGTAG
img=$DC_MQIMGREG/txmq-mq-base-rpm-$MQVER:$tag

qmenv=$1

. $qmenv

qmname=$MQ_QMGR_NAME

if [[ -z $qmname ]]; then
echo qmgr environment file required
exit 1
fi

net=qmnet

# vault
vault="-e VAULT_ENABLE_TLS=$VAULT_ENABLE_TLS -e VAULT_LDAP_CREDS_INJECT_PATH=$VAULT_LDAP_CREDS_INJECT_PATH -e VAULT_TLS_KEY_INJECT_PATH=$VAULT_TLS_KEY_INJECT_PATH -e VAULT_TLS_CERT_INJECT_PATH=$VAULT_TLS_CERT_INJECT_PATH -e VAULT_TLS_CA_INJECT_PATH=$VAULT_TLS_CA_INJECT_PATH"

# git
git="-e GIT_CONFIG_URL=$GIT_CONFIG_URL -e GIT_CONFIG_REF=$GIT_CONFIG_REF -e GIT_CONFIG_DIR=$GIT_CONFIG_DIR"

# web
web="-e MQ_START_MQWEB=$MQ_START_MQWEB"

# debug
debug="-e MQRUNNER_DEBUG=$MQRUNNER_DEBUG"

# tls
tls="-e MQ_ENABLE_TLS_NO_VAULT=$MQ_ENABLE_TLS_NO_VAULT"

# ldap
ldap="-e LDAP_BIND_PASSWORD=$LDAP_BIND_PASSWORD"

# filter mq log to standard output
# comma separated list of prefixes
# empty value will suppress mq output to std out
# special value NO_FILTER will output every line of mq log
# special value DEFAULT_FILTER will apply AMQ filter to mq output
logfilter="-e MQ_LOG_FILTER=$MQ_LOG_FILTER"

# qmgr, required
qmgr="-e MQ_QMGR_NAME=$qmname"

# all envs
envars="$qmgr $debug $tls $web $vault $git $ldap $logfilter"

# run
sudo podman run --rm --name $qmname --network $net -v mqdata:/var/mqm -v qmini:/etc/mqm/qmini -v mqsc:/etc/mqm/mqsc -v qmtls:/etc/mqm/pki/cert -v qmtrust:/etc/mqm/pki/trust -v webuser:/etc/mqm/webuser $envars -p 1414:1414 -p 9443:9443 -p 40000:40000 $img
