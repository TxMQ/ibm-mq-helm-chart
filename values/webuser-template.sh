webuser:
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
      host: openldap.default.svc.cluster.local
      port: 389
      ldaptype: Custom
      binddn: cn=admin,dc=mqldap,dc=com
      bindpassword: admin
      basedn: dc=mqldap,dc=com
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
