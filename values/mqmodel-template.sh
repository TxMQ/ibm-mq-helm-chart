#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo env file parameter required: mqmodel-template.sh \<envfile\>
exit 1
fi

# read env
. $envfile

outdir=output

cat <<EOF > $outdir/mqmodel.yaml
mq:
  qmgr:
    name: $QMNAME
    access:
      allowip: ["*"]
    authority:
    - group: [devs]
      grant: [connect, inq]
    - group: [devs]
      grant: [alladm]
    alter: []

  auth:
    ldap:
      connect:
        ldaphost: "$LDAP_HOST"
        ldapport: $LDAP_PORT
        binddn: "$LDAP_USER"
        bindpassword: ""
        tls: false
      groups:
        groupsearchbasedn: "$BASEDN_GROUPS"
        objectclass: "groupOfNames"
        groupnameattr: "cn"
        groupmembershipattr: "member"
      users:
        usersearchbasedn: "$BASEDN_USERS"
        objectclass: "inetOrgPerson"
        usernameattr: "uid"
        shortusernameattr: "cn"

  svrconn:
  - svrconnproperties:
      name: APP.SVRCONN
      maxmsgl: 4096
    tls:
      enabled: true
      clientauth: true
      ciphers: [TLS_RSA_WITH_AES_128_CBC_SHA256]
    access:
      allowip: ['*']
    authority:
      - group: [devs]
        grant: [chg, crt, dlt, dsp, ctrl, ctrlx]
      - group: [devs]
        grant: [alladm]
    alter:
      - ALTER CHANNEL(APP.SVRCONN) CHLTYPE(SVRCONN) SSLCAUTH(OPTIONAL)

  localqueue:
  - name: q.a

    defaultprioprity: 2
    defaultpersistence: true

    maxmsgl: 4096
    maxdepth: 1000

    authority:
    - group: [devs]
      grant: [put, get, dsp]
    - group: [devs]
      grant: [alladm]
      revoke: [dlt]
EOF
