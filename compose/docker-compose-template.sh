#!/bin/bash -x

. ../env.sh

ldapenv=$1
qmgrenv=$2

if [[ -z $ldapenv ]]; then
echo ldap env file required: docker-compose-template.sh '<ldapenv> <qmgrenv>'
exit 1
fi

if [[ -z $qmgrenv ]]; then
echo qmgr env file required: docker-compose-template.sh '<ldapenv> <qmgrenv>'
exit 1
fi

# load env
. $ldapenv
. $qmgrenv

outdir=output
qmname=$MQ_QMGR_NAME

cat <<EOF > $outdir/docker-compose.yaml
version: "3.9"

services:
  openldap:
    image: docker.io/bitnami/openldap:latest
    ports:
      - '$LDAP_TCP_PORT:1389'
      - '$LDAP_SSL_PORT:1636'
    volumes:
      - ldif:/ldifs
    environment:
      - LDAP_ROOT=$LDAP_ROOT
      - LDAP_ADMIN_USERNAME=$LDAP_ADMIN_USERNAME
      - LDAP_ADMIN_PASSWORD=$LDAP_ADMIN_PASSWORD
      - LDAP_ALLOW_ANON_BINDING='no'

  mqrunner:
    image: $DC_MQIMGREG/txmq-mq-base-rpm-$MQVER:$MQIMGTAG
    depends_on:
      - openldap
    ports:
      - '1414:1414'
      - '9443:9443'
      - '40000:40000'
    volumes:
      - mqmq:/mnt/data/mqm
      - mqmd:/mnt/data/md
      - mqld:/mnt/data/ld
      - qmini:/etc/mqm/qmini
      - mqsc:/etc/mqm/mqsc
      - qmtls:/etc/mqm/pki/cert 
      - qmtrust:/etc/mqm/pki/trust
      - webuser:/etc/mqm/webuser
    environment:
      - MQ_QMGR_NAME=$qmname
      - LDAP_BIND_PASSWORD=$LDAP_BIND_PASSWORD
      - MQ_START_MQWEB=1
      - MQRUNNER_DEBUG=1
      - MQ_LOG_FILTER=
      - GIT_CONFIG_URL= 
      - GIT_CONFIG_REF= 
      - GIT_CONFIG_DIR=
      - MQ_ENABLE_TLS_NO_VAULT=1
      - VAULT_ENABLE_TLS=false 
      - VAULT_LDAP_CREDS_INJECT_PATH= 
      - VAULT_TLS_KEY_INJECT_PATH= 
      - VAULT_TLS_CERT_INJECT_PATH= 
      - VAULT_TLS_CA_INJECT_PATH=

volumes:
  mqmq:
  mqmd:
  mqld:

  qmini:
    external: true

  mqsc:
    external: true

  qmtls:
    external: true

  qmtrust:
    external: true

  webuser:
    external: true

  ldif:
    external: true
EOF
