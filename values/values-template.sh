#!/bin/bash

# read env
. ../env.sh

outdir=$1
qmname=$2

cat << EOF > $outdir/values.yaml
qmspec:

  #
  # accept license - required
  #
  license:
    accept: 'true'

  #
  # create image pull secret
  # oc create secret docker-registry image-pull-secret --docker-username=<u> --docker-password=<p> --docker-email=<e>
  #
  imagePullSecrets: 
    - name: 'image-pull-secret'

  #
  # create tls secret
  # oc create secret generic qm-tls --from-file=tls.key=</path/to/tls.key> --from-file=tls.crt=</path/to/tls.crt> --from-file=ca.crt=</path/to/ca.crt>
  #
  # create trust config map with trust certificates
  # config map keys must have .crt suffix.
  # inlcude as many ca files as you need.
  # oc create configmap qm-trust --from-file=ca1.crt=</path/to/ca1.crt> --from-file=ca2.crt=</path/to/ca2.crt> ...
  #
  pki:
    tlsSecretName: 'qm-tls'
    trustMapName: 'qm-trust'

  #
  # create ldap secret
  # oc create secret generic ldapcreds --from-literal=password=<ldappassword>
  #
  ldapCredsSecret:
    name: 'ldapcreds' # ldapcreds

  # queue manager name - required
  name: $qmname

  # custom docker image - required
  image: $MQIMGREG/txmq-mq-base-rpm-$MQVER:$MQIMGTAG

  # image pull policy IfNotPresent|Always
  imagePoolPolicy: Always

  #
  # set environment variables
  #
  env:
    # start mq web console
  - name: MQ_START_MQWEB
    value: "1"
  - name: GIT_CONFIG_URL
    value: ""
  - name: GIT_CONFIG_REF
    value: ""
  - name: GIT_CONFIG_DIR
    value: ""
  - name: MQRUNNER_DEBUG
    value: "1"

  - name: MQ_LOG_FILTER
    # filter mq log to standard output
    # comma separated list of prefixes
    # empty value will suppress mq output to std out
    # special value NO_FILTER will output every line of mq log
    # special value DEFAULT_FILTER will apply AMQ filter to mq output
    #
    value: ""

#  resources:
#    limits:
#      cpu: "250m"
#      memory: "512Mi"
#    requests:
#      cpu: "250m"
#      memory: "512Mi"

  storage:
    usePvc: 'true'
    pvcName: qm-sts-claim
    storageClass: standard
    accessMode: ReadWriteOnce
    # @todo
    deleteClaim: false
    size: 2Gi

#
# mq web console configuration
#
webuser:
  #
  # groups mapped to built-in web roles
  #
  webroles:
  - name: MQWebAdmin
    groups: [devs]
  - name: MQWebAdminRO
    groups: [devs]
  - name: MQWebUser
    groups: [devs]

  #
  # groups mapped to built-in api roles
  #
  apiroles:
  - name: MQWebAdmin
    groups: [devs]
  - name: MQWebAdminRO
    groups: [devs]
  - name: MQWebUser
    groups: ["devs"]

  #
  # ldap registry
  #
  ldapregistry:
    connect:
      realm: openldap
      host: openldap.default.svc.cluster.local
      port: 389
      ldaptype: Custom
      binddn: cn=admin,dc=mqldap,dc=com
      #
      # bind password is defined in ldap secret
      # or in the vault if configured
      # otherwise define bind password here.
      #
      bindpassword: ""
      basedn: dc=mqldap,dc=com
      sslenabled: false

    # group metadata
    groupdef:
      objectclass: groupOfNames
      groupnameattr: cn
      groupmembershipattr: member

    # user metadata
    userdef:
      objectclass: inetOrgPerson
      usernameattr: uid
EOF
