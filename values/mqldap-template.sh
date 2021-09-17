#!/bin/bash

envfile=$1

if [[ -z ${envfile} ]]; then
echo env file parameter required: mqldap-template.sh 'envfile'
exit 1
fi

# read env
. ${envfile}

outdir=output

if [[ ${LDAP_TYPE} == "activedirectory" ]]; then
userobjectclass="USER"
usernameattr="sAMAccountName"
shortuser="employeeID"
groupobjectclass="GROUP"
groupnameattr="sAMAccountName"
groupmembershipattr="member"
else
# default: openldap
userobjectclass="inetOrgPerson"
usernameattr="uid"
shortuser="cn"
groupobjectclass="groupOfNames"
groupnameattr="cn"
groupmembershipattr="member"
fi

cat <<EOF > $outdir/mqldap.yaml
mq:
  qmgr:
    name: ${QMNAME}
    alter: []

  auth:
    ldap:
      connect:
        ldaphost: "${LDAP_HOST}"
        ldapport: ${LDAP_PORT}
        binddn: "${LDAP_USER}"
        bindpassword: ""
        tls: false
      groups:
        groupsearchbasedn: "${BASEDN_GROUPS}"
        objectclass: "${groupobjectclass}"
        groupnameattr: "${groupnameattr}"
        groupmembershipattr: "${groupmembershipattr}"
      users:
        usersearchbasedn: "${BASEDN_USERS}"
        objectclass: "${userobjectclass}"
        usernameattr: "${usernameattr}"
        shortusernameattr: "${shortuser}"

  alter: []
EOF
