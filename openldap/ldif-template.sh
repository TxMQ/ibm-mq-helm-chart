#!/bin/bash

# run this script to generate bitnami-ldif config map in the output directory
# you can edit generated config map to match your requirements
# apply generated config map

mkdir -p output

ldaproot=${LDAP_ROOT:-dc=mqldap,dc=com}
userpassword=${LDAP_USER_PASSWORD:-hello}

read -r -a dcarr <<<$(echo "$ldaproot" | tr "=," " ")
dc=${dcarr[1]}
o=${dcarr[1]}

out=output/ldif-config-map.yaml

if [[ -f $out ]]; then
echo Config map file \"$out\" already exists. Delete and re-run.
exit 1
fi

cat <<EOF > $out
apiVersion: v1
kind: ConfigMap
metadata:
  name: bitnami-ldif
data:
  bootstrap.ldif: |
    dn: $ldaproot
    objectClass: dcObject
    objectClass: organization
    dc: $dc
    o: $o

    dn: ou=users,$ldaproot
    objectClass: organizationalUnit
    ou: users

    dn: ou=groups,$ldaproot
    objectClass: organizationalUnit
    ou: groups

    dn: uid=user1,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user1
    cn: user1
    sn: user1
    userPassword: $userpassword

    dn: uid=user2,ou=users,dc=$ldaproot
    objectClass: inetOrgPerson
    uid: user2
    cn: user2
    sn: user2
    userPassword: $userpassword

    dn: uid=user3,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user3
    cn: user3
    sn: user3
    userPassword: $userpassword

    dn: uid=user4,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user4
    cn: user4
    sn: user4
    userPassword: $userpassword

    dn: uid=user5,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user5
    cn: user5
    sn: user5
    userPassword: $userpassword

    dn: uid=user6,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user6
    cn: user6
    sn: user6
    userPassword: $userpassword

    dn: uid=user7,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user7
    cn: user7
    sn: user7
    userPassword: $userpassword

    dn: uid=user8,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user8
    cn: user8
    sn: user8
    userPassword: $userpassword

    dn: uid=user9,ou=users,$ldaproot
    objectClass: inetOrgPerson
    uid: user9
    cn: user9
    sn: user9
    userPassword: $userpassword

    dn: cn=devs,ou=groups,$ldaproot
    objectClass: groupOfNames
    cn: devs
    member: uid=user1,ou=users,$ldaproot
    member: uid=user2,ou=users,$ldaproot
    member: uid=user3,ou=users,$ldaproot
    member: uid=user4,ou=users,$ldaproot

    dn: cn=admins,ou=groups,$ldaproot
    objectClass: groupOfNames
    cn: admins
    member: uid=user5,ou=users,$ldaproot
    member: uid=user6,ou=users,$ldaproot
    member: uid=user7,ou=users,$ldaproot

    dn: cn=readadmins,ou=groups,$ldaproot
    objectClass: groupOfNames
    cn: readadmins
    member: uid=user8,ou=users,$ldaproot
    member: uid=user9,ou=users,$ldaproot
EOF
