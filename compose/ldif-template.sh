#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo env file required: ldif-template.sh '<envfile>'
exit 1
fi

# load env
. $envfile

# compute dc and o from ldap_root
dc=mqldap
o=mqldap

out=output/ldif

cat <<EOF > $out/bootstrap.ldif
dn: $LDAP_ROOT
objectClass: dcObject
objectClass: organization
dc: $dc
o: $o

dn: ou=users,$LDAP_ROOT
objectClass: organizationalUnit
ou: users

dn: ou=groups,$LDAP_ROOT
objectClass: organizationalUnit
ou: groups

dn: uid=user1,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user1
cn: user1
sn: user1
userPassword: hello

dn: uid=user2,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user2
cn: user2
sn: user2
userPassword: hello

dn: uid=user3,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user3
cn: user3
sn: user3
userPassword: hello

dn: cn=devs,ou=groups,$LDAP_ROOT
objectClass: groupOfNames
cn: devs
member: uid=user1,ou=users,$LDAP_ROOT
member: uid=user2,ou=users,$LDAP_ROOT
member: uid=user3,ou=users,$LDAP_ROOT
EOF
