### Builiding custom IBM MQ container image.

Custom image is build from the rpm distribution of IBM MQ.

IBM MQ rpm distribution is not publicly available and must be downloaded
from IBM Passport Advantage to the local directory.<br>

Clone build git repository to your local machine.<br>

Each release is based on specific MQ version. MQ version is compiled into the chart release.<br>

**Build custom image.**

```
git clone ...
cd mq-operator && mkdir rpm
tar xvf ... -C rpm

sudo podman login docker.io
build.sh docker-repo`
```

**Create image pull secret**

`oc create secret docker-registry image-pull-secret --docker-username=<u> --docker-password=<p> --docker-email=<e>`

### TxMQ MQ helm chart.

To make it easier to work with the TxMQ chart configuration, yaml objects are grouped
into a number of files and are available in the `values` directory.<br>

- qmgr.yaml
- vault.yaml
- qmini.yaml
- mqscic.yaml
- mq.yaml

*qmgr.yaml* defines basic kubernetes and queue manager configuration.<br>

All other files are optional.<br>

Use `helm` command to install TxMQ chart.<br>

`helm install -f qmgr.yaml [-f vault.yaml] [-f mqscic.yaml] [-f qmini.yaml] [-f mq.yaml] release mqchart/`

**Dependencies**<br>
Ldap server is required.<br>
You can either use existing LDAP server or deploy openldap chart.<br>

Hashicorp vault is recommended. You can either use existing vault, or deploy hashicorp vault chart.<br>

**Configuring a chart**<br>
Chart configuration contains many parameters and and be updated at any time.<br>

`values` directory contains starter configuration files.

Look at the comments and update `values/values.yaml` file for the first configuration.<br>

That would include accepting a license, creating kubernetes secrets, naming queue manager, and ldap configuration.<br>

`helm install -f values/values.yaml <release-name> mqchart/`

It is recommended that Hashicorp vault integration is configured as next step.<br>

After vault is configured in `values/vault.yaml` pass it to the chart:<br>
`helm install -f values/values.yaml -f values/vault.yaml <release-name> mqchart/`

To configure queue manager at startup, place *mqsc* commands in `values/mqscic.yaml` file and pass it to the chart:
`helm install -f values/values.yaml -f values/vault.yaml -f values/mqscic.yaml <release-name> mqchart/`

To update queue manager ini parameters, place them in the `values/qmini.yaml` file and pass it to the chart:
`helm install -f values/values.yaml -f values/vault.yaml -f values/mqscic.yaml -f values/qmini.yaml <release-name> mqchart/`

To use higher-level mq configuration abstraction, put it in the `values/mq.yaml` file and pass it to the chart:<br>
`helm install -f values/values.yaml -f values/vault.yaml -f values/mqscic.yaml -f values/qmini.yaml -f values/mq.yaml <release-name> mqchart/`

The first file is required. Any combination of files can be passed to the chart.<br>

** Mq git configuration **<br>

**Examples and Reference**<br>

Use examples as a starting point for chart configuration.

**Secrets.**<br>
Secrets are used for LDAP authentication and TLS keys and certificates.

You can configure secrets with or without a vault.

Vault configuration takes precedence over kubernetes secrets.

**Kubernetes LDAP secrets.**<br>
When secret vault is not used create generic kubernetes secret with the *password* key.

`oc create secret generic qm-ldap-creds --from-literal=password=ldappassword`

Set LDAP secret name in yaml configuration:

```
qmspec:
  ldapCredsSecret:
    name: qm-ldap-creds
```

#### Kubernetes TLS secrets.
When secret vault is not used create generic kubernetes secret with the *tls.key*, *tls.crt*, and *ca.crt* keys.

`oc create secret generic qm-tls --from-file=tls.key=/path/to/tls.key --from-file=tls.crt=/path/to/tls.crt --from-file=ca.crt=/path/to/ca.crt`

Set TLS secret name in yaml configuration:
```
qmspec:
  pki:
    tlsSecretName: qm-tls
```

### Integration with Hashicorp vault.

#### Vault prerequisites.

**Install hashicorp vault helm chart.**

```
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm install vault hashicorp/vault [--set "global.openshift=true"] --set "server.dev.enabled=true"
```

**Configure kubernetes authentication for the vault:**
```
oc exec -it vault-0 -- /bin/sh
/ $ vault auth enable kubernetes
Success! Enabled kubernetes auth method at: kubernetes/
```

**Configure kubernetes auth to use service account token.**
```
/ $ vault write auth/kubernetes/config \
token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) \
kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
Success! Data written to: auth/kubernetes/config
```

### Vault Secrets and Access Policies.

**Create LDAP vault secret.**

For ldap we create a secret path with one key-value pair: password=value

```
/ $ vault kv put secret/mq/ldapcreds password="admin"
Key              Value
---              -----
created_time     2021-07-22T21:12:41.138756172Z
deletion_time    n/a
destroyed        false
version          1
```

**Create TLS vault secret.**

For tls we create a secret path with 3 key/value pairs:
private key (key.pem), cert (cert.pem), and ca chain (ca.pem)

Each key/value will be injected into separate path.

Copy key.pem, cert.pem, and ca.pem files to the vault container:

```
oc cp ./tls.key vault-0:/home/vault
oc cp ./tls.crt vault-0:/home/vaul
```

**Create TLS secret with 3 key-value pairs: key, cert, and ca**

`vault kv put secret/mq/tls key=@tls.key cert=@lts.crt ca=@ca.crt`

**Vault Access Control.**

Define policy to allow read access to ldap creds and tls secrets.

Path values in the policy are derived from the secret paths but not the same.
Note that 'data' segment is injected into the secret path.

```
vault policy write mq - <<EOF
path "secret/data/mq/ldapcreds" {
    capabilities = ["read"]
}

path "secret/data/mq/tls" {
    capabilities = ["read"]
}
EOF
```

**Bind service account and namespace to a policy to create a role:**

Vault authorizes specific service account to connet and get secret token.
Service account is created at chart startup time and is prefixed with the name of the helm release.

Suppose that chart release `bar` is deployed in namespace `foo`.
Then service account name is `foo-mqdeployer`

Create kubernetes authentication role by binding policy to service account and namespace.

```
vault write auth/kubernetes/role/mq \
bound_service_account_names=foo-mqdeployer \
bound_service_account_namespaces=bar \
policies=mq
```

### Vault secret injection annotations.

Secrets from the vault are injected by the vault agent injector that is deployed by the vault helm chart.

Injection is driven by annotations applied to the queue manager pod.

Annotation explanation.

Enable secret injection:
`vault.hashicorp.com/agent-inject: "true"`

Vault authentication role:
`vault.hashicorp.com/role: mq`

Inject ldap vault secret into /vault/secrets/mq-ldapcreds.txt file
`vault.hashicorp.com/agent-inject-secret-[mq-ldapcreds.txt]: secret/data/mq/ldapcreds`

Set Ldap vault secret file template.
```
vault.hashicorp.com/agent-inject-template-[mq-ldapcreds.txt]: |
  {{- with secret "secret/data/mq/ldapcreds" -}}
  {{ .Data.data.password }}
  {{- end -}}
```

Inject tls vault secret into /vault/secrets/tls.key
`vault.hashicorp.com/agent-inject-secret-[tls.key]: secret/data/mq/tls`

Tls key secret template:
```
vault.hashicorp.com/agent-inject-template-[tls.key]: |
  {{- with secret "secret/data/mq/tls" -}}
  {{ .Data.data.key }}
  {{- end -}}
```

Inject tls vault secret into /vault/secrets/tls.crt
`vault.hashicorp.com/agent-inject-secret-[tls.crt]: secret/data/mq/tls`

Tls cert secret template:
```
vault.hashicorp.com/agent-inject-template-[tls.crt]: |
  {{- with secret "secret/data/mq/tls" -}}
  {{ .Data.data.cert }}
  {{- end -}}
```  

**vault injection annotations example**
```
qmspec:
  annotations:
    vault.hashicorp.com/agent-inject: 'true'

    # role for ldap creds and tls key pair
    vault.hashicorp.com/role: 'mq'

    # ldap creds
    vault.hashicorp.com/agent-inject-secret-mq-ldapcreds.txt: 'secret/data/mq/ldapcreds'
    vault.hashicorp.com/agent-inject-template-mq-ldapcreds.txt: |          
      {{- with secret "secret/data/mq/ldapcreds" -}}
      {{ .Data.data.password }}
      {{- end -}}

    # tls key pair
    vault.hashicorp.com/agent-inject-secret-tls.key: 'secret/data/mq/tls'
    vault.hashicorp.com/agent-inject-template-tls.key : |
      {{- with secret "secret/data/mq/tls" -}}
      {{ .Data.data.key }}
      {{- end -}}
    vault.hashicorp.com/agent-inject-secret-tls.crt: 'secret/data/mq/tls'
    vault.hashicorp.com/agent-inject-template-tls.crt : |
      {{- with secret "secret/data/mq/tls" -}}
      {{ .Data.data.cert }}
      {{- end -}}
    vault.hashicorp.com/agent-inject-secret-ca.crt: 'secret/data/mq/tls'
    vault.hashicorp.com/agent-inject-template-ca.crt : |
      {{- with secret "secret/data/mq/tls" -}}
      {{ .Data.data.ca }}
      {{- end -}}
```

**Enable vaut in mq chart.**

```
qmspec:
  vault:
    ldapCreds:
      enable: 'true'
      injectpath: '/vault/secrets/ldapcreds.txt'
    tls:
      enable: 'true'
      keyinjectpath: '/vault/secrets/tls.key'
      certinjectpath: '/vault/secrets/tls.crt'
      cainjectpath: '/vault/secrets/ca.pem'
```

### Configuring queue manager kubernetes parameters.

```
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
```

# Configuring mq web console.

Mq web console requires definition of a number of predefined roles,
authentication, and key store.

Groups that are defined in webuser must be authorized to access queue manager in queue manager configuration.

If you want to enable mqwebconsole it must be configured to
authenticate against the same ldap server as queue manager.

## Configuring mq web console authentication.

```
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
```

## Configuring mq web console roles.

```
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
```

## Configuring queue manager authentication.

Queue manager running in kubernetes cluster must be configured to
authenticate against ldap server.

Queue manager ldap authentication can be configured as mqsc command
to be run at queue manager startup or by using higher level txmq abstraction
to configure ldap authentication in yaml format in mq.yaml file. If mq.yaml
file is used it is transformed to mqsc commands before queue manager startup
and merged with other native mqsc startup commands.

Here we show high-level yaml configuration for ldap server.

```
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
```

### integration with persistent storage

### TxMQ Chart Reference.

**qmspec** object:<br>

| Path                        | Type       | Value    |
| :---                        | :---:      | :---     |
|license.accept               | string     | 'true'   |
|labels                       | map        | additional labels for chart resources |
|annotations                  | map        | additional annotations for chart resources |
|capabilities                 | string     | mq image capabilities: 'mqbase' |
|licenseAnnotations           | map        | custom license annotations to apply to chart resources |
|affinity                     | map        | kubernetes affinity object for queue manager pods |
|serviceAccount               | map        | service account object
|serviceAccount.name          | string     | service account name, default 'mqdepolyer' |
|serviceAccount.create        | boolean    | true - chart will create service account, false - existing service account |
|imagePullSecrets             | map        | image pull secrets object
|imagePullSecrets.name        | string     | docker registry image pull secret |
|pki                          | map        | queue manager pki object
|pki.tlsSecretName            | string     | generic TLS secret. Not recommended, use vault instead |
|ldapCredsSecret              |            | queue manager LDAP credentials object
|ldapCredsSecret.name         |string      | generic TLS secret name. Not recommended, use vault instead |
|ldapCredsSecret.passwordKey  |string      | secret key, default: password
|vault                        | object     |hashicorp vault object|
|vault.ldapCreds              | object     |vault LDAP credentials object |
|vault.ldapCreds.enable       | string      | 'true' - inject queue manager LDAP credentials from the vault, defaults to 'false' |
|vault.ldapCreds.injectpath   | string      | vault credentials injection path. Prefix with /vault/secrets/
|vault.tls                    | object      |vault TLS credentials object|
|vault.tls.enable             | string      | 'true' - inject queue manager TLS credentials from the vault, defaults to 'false' |
|vault.tls.keyinjectpath      | string      | TLS key injection path, prefix with /vault/secret |
|vault.tls.certinjectpath     | string      | TLS cert injection path, prefix with /vault/secret |
|vault.tls.cainjectpath       | string      | TLS ca injection path, prefix with /vault/secret |
|terminationGracePeriodSeconds | integer    | 10 seconds. @todo |
|storage                       | object     | storage object. @todo |
|storage.usepvc                | string     | 'true' - use pvc for persistent storage |
|pvcname                       | string     | chart external pvc name |
