#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo env file required: ./webuser-template.sh \<envfile\>
exit 1
fi

# load env
. $envfile

if [[ $LDAP_TYPE == "activedirectory" ]]; then
# active directory
groupdef="
    groupdef:
      objectclass: GROUP
      groupnameattr: sAMAccountName
      groupmembershipattr: member
"
userdef="
    userdef:
      objectclass: USER
      usernameattr: sAMAccountName
"
else
# default: openldap
groupdef="
    groupdef:
      objectclass: groupOfNames
      groupnameattr: cn
      groupmembershipattr: member
"
userdef="
    userdef:
      objectclass: inetOrgPerson
      usernameattr: uid
"
fi

cat <<EOF > output/webuser.yaml
webuser:
  webroles:
  - name: MQWebAdmin
    groups: ["$ADMIN_GROUP"]
  - name: MQWebAdminRO
    groups: ["$READ_ADMIN_GROUP"]
  - name: MQWebUser
    groups: ["$APPL_GROUP"]

  apiroles:
  - name: MQWebAdmin
    groups: ["$ADMIN_GROUP"]
  - name: MQWebAdminRO
    groups: ["$READ_ADMIN_GROUP"]
  - name: MQWebUser
    groups: ["$APPL_GROUP"]

  ldapregistry:
    connect:
      realm: openldap
      host: $LDAP_HOST
      port: $LDAP_PORT
      ldaptype: Custom
      binddn: $LDAP_USER
      bindpassword: 
      basedn: $LDAP_ROOT
      sslenabled: false
$groupdef
$userdef
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
