to add ldif file entries:

ldapadd -x -D "cn=admin,dc=mqldap,dc=com" -w admin -H ldap:// -f ./bootstrap.ldif

to search:

ldapsearch -x -D "cn=admin,dc=mqldap,dc=com" -w admin -H ldap:// -b "dc=mqldap,dc=com"

ldapsearch -x -H ldap://openldap.default.svc.cluster.local:389 -b "dc=mqldap,dc=com" -D "cn=admin,dc=mqldap,dc=com" -w admin
