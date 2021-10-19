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
    dn: dc=mqldap,dc=com
    objectClass: dcObject
    objectClass: organization
    dc: mqldap
    o: mqldap

    dn: ou=users,dc=mqldap,dc=com
    objectClass: organizationalUnit
    ou: users

    dn: ou=groups,dc=mqldap,dc=com
    objectClass: organizationalUnit
    ou: groups

    dn: uid=user1,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user1
    cn: user1
    sn: user1
    userPassword: hello

    dn: uid=user2,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user2
    cn: user2
    sn: user2
    userPassword: hello

    dn: uid=user3,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user3
    cn: user3
    sn: user3
    userPassword: hello

    dn: uid=user4,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user4
    cn: user4
    sn: user4
    userPassword: hello

    dn: uid=user5,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user5
    cn: user5
    sn: user5
    userPassword: hello

    dn: uid=user6,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user6
    cn: user6
    sn: user6
    userPassword: hello

    dn: uid=user7,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user7
    cn: user7
    sn: user7
    userPassword: hello

    dn: uid=user8,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user8
    cn: user8
    sn: user8
    userPassword: hello

    dn: uid=user9,ou=users,dc=mqldap,dc=com
    objectClass: inetOrgPerson
    uid: user9
    cn: user9
    sn: user9
    userPassword: hello

    dn: cn=devs,ou=groups,dc=mqldap,dc=com
    objectClass: groupOfNames
    cn: devs
    member: uid=user1,ou=users,dc=mqldap,dc=com
    member: uid=user2,ou=users,dc=mqldap,dc=com
    member: uid=user3,ou=users,dc=mqldap,dc=com
    member: uid=user4,ou=users,dc=mqldap,dc=com

    dn: cn=admins,ou=groups,dc=mqldap,dc=com
    objectClass: groupOfNames
    cn: devs
    member: uid=user5,ou=users,dc=mqldap,dc=com
    member: uid=user6,ou=users,dc=mqldap,dc=com
    member: uid=user7,ou=users,dc=mqldap,dc=com

    dn: cn=readadmins,ou=groups,dc=mqldap,dc=com
    objectClass: groupOfNames
    cn: devs
    member: uid=user8,ou=users,dc=mqldap,dc=com
    member: uid=user9,ou=users,dc=mqldap,dc=com
EOF
