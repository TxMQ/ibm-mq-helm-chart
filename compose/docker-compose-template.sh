#!/bin/bash -x

. ../env.sh

outdir=$1
qmgrname=$2

cat <<EOF > $outdir/docker-compose.yaml
version: "3.9"

services:
  openldap:
    image: docker.io/bitnami/openldap:latest
    ports:
      - '1389:1389'
      - '1636:1636'
    volumes:
      - ldif:/ldifs
    environment:
      - LDAP_ROOT=dc=mqldap,dc=com
      - LDAP_ADMIN_USERNAME=admin
      - LDAP_ADMIN_PASSWORD=admin
      - LDAP_ALLOW_ANON_BINDING=no

  mqrunner:
    image: $DC_MQIMGREG/txmq-mq-base-rpm-$MQVER:$MQIMGTAG
    depends_on:
      - openldap
    ports:
      - '1414:1414'
      - '9443:9443'
      - '40000:40000'
    volumes:
      - mqdata:/var/mqm
      - mqsc:/etc/mqm/mqsc
      - qmtls:/etc/mqm/pki/cert 
      - qmtrust:/etc/mqm/pki/trust
      - webuser:/etc/mqm/webuser
    environment:
      - MQ_QMGR_NAME=$qmgrname
      - LDAP_BIND_PASSWORD=admin
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
  mqdata:

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
