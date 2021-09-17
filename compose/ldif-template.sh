#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo env file required: ldif-template.sh '<envfile>'
exit 1
fi

# load env
. $envfile

# compute dc and o from ldap_root
# we assume 2 valued ldap root, eg: dc=mqldap,dc=com
read -r -a dcarr <<<$(echo "$LDAP_ROOT" | tr "=," " ")
dc=${dcarr[1]}
o=${dcarr[1]}

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

dn: uid=user4,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user4
cn: user4
sn: user4
userPassword: hello

dn: uid=user5,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user5
cn: user5
sn: user5
userPassword: hello

dn: uid=user6,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user6
cn: user6
sn: user6
userPassword: hello

dn: uid=user7,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user7
cn: user7
sn: user7
userPassword: hello

dn: uid=user8,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user8
cn: user8
sn: user8
userPassword: hello

dn: uid=user9,ou=users,$LDAP_ROOT
objectClass: inetOrgPerson
uid: user9
cn: user9
sn: user9
userPassword: hello

dn: cn=devs,ou=groups,$LDAP_ROOT
objectClass: groupOfNames
cn: devs
member: uid=user1,ou=users,$LDAP_ROOT
member: uid=user2,ou=users,$LDAP_ROOT
member: uid=user3,ou=users,$LDAP_ROOT
member: uid=user4,ou=users,$LDAP_ROOT

dn: cn=admins,ou=groups,$LDAP_ROOT
objectClass: groupOfNames
cn: admins
member: uid=user5,ou=users,$LDAP_ROOT
member: uid=user6,ou=users,$LDAP_ROOT
member: uid=user7,ou=users,$LDAP_ROOT

dn: cn=readadmins,ou=groups,$LDAP_ROOT
objectClass: groupOfNames
cn: readadmins
member: uid=user8,ou=users,$LDAP_ROOT
member: uid=user9,ou=users,$LDAP_ROOT
EOF
