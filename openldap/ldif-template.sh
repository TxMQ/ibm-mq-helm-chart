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

    dn: cn=devs,ou=groups,dc=mqldap,dc=com
    objectClass: groupOfNames
    cn: devs
    member: uid=user1,ou=users,dc=mqldap,dc=com
    member: uid=user2,ou=users,dc=mqldap,dc=com
    member: uid=user3,ou=users,dc=mqldap,dc=com