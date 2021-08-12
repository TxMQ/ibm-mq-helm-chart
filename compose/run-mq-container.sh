#!/bin/bash -x

# vault
vault="-e VAULT_ENABLE_TLS=false -e VAULT_LDAP_CREDS_INJECT_PATH= -e VAULT_TLS_KEY_INJECT_PATH= -e VAULT_TLS_CERT_INJECT_PATH= -e VAULT_TLS_CA_INJECT_PATH="

# git
git="-e GIT_CONFIG_URL= -e GIT_CONFIG_REF= -e GIT_CONFIG_DIR="

# web
web="-e MQ_START_WEB=0"

# debug
debug="-e MQRUNNER_DEBUG=1"

# qmgr, required
qmgr="-e MQ_QMGR_NAME=qm20"

# all envs
envars="$qmgr $debug $web $vault $git"

# run
sudo podman run --rm -v mqdata:/var/mqm $envars -p 1414:1414 localhost/txmq-mq-base-rpm-9.2.2.0:159
