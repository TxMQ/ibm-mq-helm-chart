# Builiding custom IBM MQ image

Custom image is build from the rpm distribution of IBM MQ.

IBM MQ rpm distribution is not publicly available and must be downloaded
from IBM Passport Advantage to the local directory.<br>

Clone build git repository to your local machine.<br>

Each release is based on specific mq version.
Mq version is compiled into chart release.

Build custom image.

`git clone ...`
`cd mq-operator && mkdir rpm`
`tar xvf ... -C rpm`

`sudo podman login docker.io`
`build.sh docker-repo`

After 'base' custom image is built it is possible to build
other custom images from that.

# Txmq mq helm chart.

At present, TxMQ mq chart deploys standalone queue manager.

There are a number of files that are passed as input to the helm chart:
values.yaml, mqscic.yaml, qmini.yaml, and mq.yaml.

values.yaml file specifies basic kubernetes settings for the queue manager under `.Values.qmspec object.`
It also contains `.Values.webuser` object to configure mq web console service.

You can pass mqsc commands to be executed at queue manager startup by placing it 
in the `mqscic.yaml` file under `.Values.mqscic` object.

You can pass qmini parameters to configure queue manager in the qmini.yaml file
by configuring `.Values.qmini` object.

Txmq chart defines an abstraction over mq configuration that you can place in mq.yaml
file under `.Values.mq` object.

To install txmq mq chart use helm:

`helm install -f values.yaml [-f mqscic.yaml] [-f qmini.yaml] [-f mq.yaml] release mqchart/`

The only required file is `values.yaml`.

### Configuring queue manager kubernetes parameters.

`
qmspec:
  license:
    accept: true # required to accept license

  capabilities: mqbase # mq image capabilities

  imagePullSecrets: # docker registry secret to pull queue manager image
  - name: image-pull-secret

  pki:
    tlsSecretName: mq-tls # tls secret name for queue manager key pair
    enableTls: 1 # enable tls for queue manager communication

  qmconf:
    name: qm20 # queue manager name
    image: registry/namespace/txmq-mq-base-rpm-9.2.2.0:1.0 # queue manager image

    env: # environment variables to pass to the queue manager
    - name: MQ_START_MQWEB
      value: "1"            # start mq web console

    resources: # limits on resource consumption, tracked by license service
      limits:
        cpu: '1'
        memory: 1Gi
      requests:
        cpu: '1'
        memory: 1Gi

    storage: # queue manager storage
      pvcName: qm-sts-claim
      storageClass: standard # storage class
      accessMode: ReadWriteOnce
      deleteClaim: false
      size: 2Gi
`

# Configuring mq web console.

Mq web console requires definition of a number of predefined roles,
authentication, and key store.

Groups that are defined in webuser must be authorized to access queue manager in queue manager configuration.

If you want to enable mqwebconsole it must be configured to
authenticate against the same ldap server as queue manager.

## Configuring mq web console authentication.

`
webuser:

  ldapregistry:
    connect:
      realm: openldap
      host: openldap.default.svc.cluster.local
      port: 389
      ldaptype: Custom
      binddn: cn=admin,dc=example,dc=com
      bindpassword: admin
      basedn: dc=example,dc=com
      sslenabled: false

    groupdef:
      objectclass: groupOfNames
      groupnameattr: cn
      groupmembershipattr: member

    userdef:
      objectclass: inetOrgPerson
      usernameattr: uid
`

## Configuring mq web console roles.

`
webuser:

  webroles:
  - name: MQWebAdmin
    groups: []
  - name: MQWebAdminRO
    groups: []
  - name: MQWebUser
    groups: []

  apiroles:
  - name: MQWebAdmin
    groups: []
  - name: MQWebAdminRO
    groups: []
  - name: MQWebUser
    groups: []
`

## Configuring queue manager authentication.

Queue manager running in kubernetes cluster must be configured to
authenticate against ldap server.

Queue manager ldap authentication can be configured as mqsc command
to be run at queue manager startup or by using higher level txmq abstraction
to configure ldap authentication in yaml format in mq.yaml file. If mq.yaml
file is used it is transformed to mqsc commands before queue manager startup
and merged with other native mqsc startup commands.

Here we show high-level yaml configuration for ldap server.

`
mq:
  auth:
    ldap:
      connect:
        ldaphost: "openldap.default.svc.cluster.local"
        ldapport: 389
        binddn: "cn=admin,dc=example,dc=com"
        bindpasswordsecret: "admin"
        tls: false
      groups:
        groupsearchbasedn: "ou=groups,dc=example,dc=com"
        objectclass: "groupOfNames"
        groupnameattr: "cn"
        groupmembershipattr: "member"
      users:
        usersearchbasedn: "ou=users,dc=example,dc=com"
        objectclass: "inetOrgPerson"
        usernameattr: "uid"
        shortusernameattr: "cn"
`

## ldap password

injecting ldap password secret as environment variable.

oc create secret generic ldapcreds --from-literal=password=password

set secret name:
`
qmspec:
  ldapCredsSecret:
    name: ldapcreds
    passwordkey: password # optional, default password key is "password"
`

# tls key pair

We assume that key pair and all certs are availabe.

Private key, cert, and ca cert can be created as generic k8s secret. Secret is then mounted on the volume path /secrets/pki
Secret name is set as environment variable

Tls secret can be written into the hashicorp vault. Secret is then injected by the vault agent into the pod.
Secret path is set as environment variable.

## integration with hashicorp vault.

There are a number of secrets: tls secret, ldap authentication secret, image pull secret, etc.

deploy and configure hashicorp vault, define secret paths, define secrets, define roles for reading secrets
and grant this role to the service account for the chart. Append chart release name to the service account.
Inject secrets into qmgr container and mount secrets to the /secrets path.

For tls secrets, set environment variable to point to the tls secret path. If cert path environment variable is set
it overrides default cert path directory.

For ldap secrets, set env variable to point to secret files. At tranformation time if secret path environment variable
is set read secret file and set password.

## integration with persistent storage
