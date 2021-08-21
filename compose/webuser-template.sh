#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo env file required: ./webuser-template.sh '<envfile>'
exit 1
fi

# load env
. $envfile

outdir=output/etc/mqm/webuser

cat <<EOF > $outdir/webuser.yaml
webroles:
- name: MQWebAdmin
  groups: [devs]
- name: MQWebAdminRO
  groups: [devs]
- name: MQWebUser
  groups: [devs]

apiroles:
- name: MQWebAdmin
  groups: [devs]
- name: MQWebAdminRO
  groups: [devs]
- name: MQWebUser
  groups: ["devs"]

ldapregistry:
  connect:
    realm: openldap
    host: $LDAP_HOST
    port: $LDAP_PORT
    ldaptype: Custom
    binddn: $LDAP_USER
    bindpassword: $LDAP_BIND_PASSWORD
    basedn: $LDAP_ROOT
    sslenabled: false

  groupdef:
    objectclass: groupOfNames
    groupnameattr: cn
    groupmembershipattr: member

  userdef:
    objectclass: inetOrgPerson
    usernameattr: uid

allowedhosts: []

clientauth:
  keystorepass: ""
  truststorepass: ""
  enabled: false

variables:
- name: httpsPort
  value: "9443"
- name: httpHost
  value: '*'
- name: mqRestCorsAllowedOrigints
  value: '*'
EOF
