### Builiding custom IBM MQ container image.

Custom image is build from the rpm distribution of IBM MQ.

IBM MQ rpm distribution is not publicly available and must be downloaded
from IBM Passport Advantage to the local directory.<br>

Clone build git repository to your local machine.<br>

Each release is based on specific MQ version. MQ version is compiled into the chart release.<br>

**Prerequisites**

The top level `prereq` directory contains prerequisite scripts that you can run to install podman, and configure podman for docker-compose. There are also files for nginx ingress controller and nginx ingress controller helm chart customization file.

**Build custom MQ image.**

Custom image build is driven by environment variables defined in the `env.sh` file.

```
export MQVER="9.2.3.0"

export RPMDIR=rpm923/MQServer

export MQIMGREG=$MQIMGREG

export DC_MQIMGREG=localhost
```

Obtain MQ rpm archive, eg `IBM_MQ_9.2.3_LINUX_X86-64.tar.gz`.

Match MQVER environment variable to the MQ version, eg: `MQVER=9.2.3.0`

`RPMDIR` variable is a directory for the expanded MQ tar archive. Note that MQ rpm archive contains `MQServer` directory, so create just the first path segment of the `$RPMDIR` environment variable directory.<br>.

Expand MQ rpm archive:<br>

`mkdir rpm923; tar xzvf IBM_MQ_9.2.3_LINUX_X86-64.tar.gz -C rpm923`

`MQIMGREG` environment variable is docker image registry for mq image. Login into  the`$MQIMGREG` docker registry:

`podman login $MQIMGREG`

`DC_MQIMGREG` environment variable is docker registry for MQ image for use by podman scripts and docker-compose. If we build image locally we can reuse local image for podman scripts and docker-compose. You can leave this value as is or set it to the value of `$MQIMGREG` environment variable.

Run the build script: `./build.sh`<br>

The build script will build MQ custom image, tag it and push it to the docker registry `$MQIMGREG`.<br>

You can review images with: `sudo podman images`.

```
REPOSITORY                                     TAG         IMAGE ID      CREATED       SIZE
$MQIMGREG/txmq-mq-base-rpm-9.2.3.0             206         2ba332a5d22d  5 days ago    2.53 GB
localhost/txmq-mq-base-rpm-9.2.3.0             206         2ba332a5d22d  5 days ago    2.53 GB
```

**TLS certificates**

TLS is required for MQ image. We use `tls-gen` project to generate mq container certificates and ca certificate. Mq container server and client certificates are signed by the ca key. If you have certificates from real ca, use them instead.<br>

To clone `tls-gen` project:<br>
`git clone https://github.com/michaelklishin/tls-gen`

There are 2 use cases for certificate generation.

One is 'local' use case, where mq container runs with podman script or docker-compose on alocal machine. In this case server CN is the name of a 'local' machine.

Another is 'kubernetes' use case, when mq container runs on a kubernetes cluster. In this case, server CN is ingress endpoint hostname, and alt subject name is mq service name in the kubernetes cluster.

*Generate TLS certificates for local podman scripts or docker-compose*
Change to the `tls-gen/basic` directory.

Run `make PASSWORD=password`<br>

This will generate `result` directory with all crypto material. Note that server private key is encrypted with the password parameter.<br>

To decrypt server private key:<br>
`cd result; openssl rsa -in server_key.pem -out server_key_nopass.pem`<br>

Note that `server_key_nopass.pem` is used later on by scripts loading queue manager container keystores.<br>

Make a copy of a `result` directory:<br>
`mv tls-gen/basic/result tls-gen/basic/result-local`.

Set environment variable to point to this directory:<br>
`export TLS_GEN_RESULT=/path/to/tls-gen/basic/result-local`  

*Generate TLS certificates for kubernetes deployment*  



**Running custom MQ image with podman scripts or docker-compose.**

`compose` directory contains scripts to prepare and start using MQ image either with podman scripts or with docker-compose.

Note that ldap directory is required for the queue manager container. When running with podman scripts or docker-compose bitnami ldap image is automatically pulled from docker io. You can configure ldap container with the bootstrap ldif file.

Change to the `compose` directory: `cd compose`.

If you did not set the `TLS_GEN_RESULT` environment variable as described above in the TLS section, you will need to run `copy-certs.sh` script as a separate step.

Decide on the queue manager name and run the build script. Pass the queue manager name to the build script. Default queue manager name is `qm1`.

`./build qm1`

The build script creates `output` directory with the number of files that you can edit.

```
ls -l output
-rw-rw-r--. docker-compose.yaml
drwxrwxr-x. etc
-rw-rw-r--. ldap.env
drwxrwxr-x. ldif
-rw-rw-r--. qm1.env
```

Ldap container configuration:<br>
`ldap.env` is ldap container environment file.
`ldif/bootstrap.ldif` is ldif bootstrap file loaded by the ldap container.

Queue manager container configuration:<br>
`qm1.env` is queue manager container environment file named after the queue manager.
`etc/mqm/webuser/webuser.yaml` is mq web console configuration file for mq container.
`etc/mqm/mqini` is a directory with mq ini file passed to mq container
`etc/mqm/mqsc` is a directory with mqsc files passed to mq container.
`etc/mqm/pki/cert` is a directory with key and certs passed to mq container to build key stores.
`etc/mqm/pki/trust` is a directory with ca certs passed to mq container to build trust store.<br>

Note that generated files are consistent with respect to queue manager environment file and ldap container environment file.<br>

If you did not set `TLS_GEN_RESULT` environment variable, copy certificates to the output directory:
`./copy-certs.sh path/to/tls-gen/basic/result-local`

Prepare volumes for queue manager container and ldap container:<br>
`./prepare-volumes.sh`

*Run ldap container with the podman script.*
The script takes ldap environment file as argument.  
`./run-ldap-container.sh output/ldap.env`

*Run mq container with the podman script.*
The script take queue manager env file as argument.  
`./run-mq-container.sh output/qm1.env`

*Run ldap container with docker-compose*  
`cd output; sudo docker-compose up openldap`

*Run qmgr container with docker compose*  
`cd output; sudo docker-compose up mqrunner`

### TxMQ MQ helm chart.

**MQ helm chart dependencies**

**Bitnami LDAP server image**
LDAP server is required for queue manager authentication.

You can use bitnami ldap container deployment included in the top-level `openldap` directory.  
Change to the `openldap` directory.  

Before you create ldap deployment, run `ldif-template.sh` script to create bootstrap ldif config map.  
You can customize bootstrap ldif data to match your requirements.  

Apply yaml files:  
```
kubectl apply -f output/ldap-config-map.yaml
kubectl apply -f bitnami-sa.yaml
kubectl apply -f bitnami-service.yaml
kubectl apply -f bitnami-deployment.yaml
```

You can run `apply-yaml.sh` script to apply above steps.

**Hashi Corp Vault**
Vault integration feature is experimental.  

**Building MQ helm chart**

MQ helm chart is located in the top-level git directory `mqchart` .  

Since mq version can change with mq releases, build script in the `mqchart` directory creates `output` directory and reconsiles mq chart with correct mq version.  

`cd mqchart; ./build.sh`

All chart files will be copied to the `output` directory. When interacting with the chart point to the `output` directory location.  

To use the chart we need to pass customization file (or files) to the `helm install` command.  

Top-level `values` directory contains build script and templates to generate values files in the `output` directory. You can make changes to the values files in the `output` directory.  

The build script in the `values` directory requires queue manager name parameter (qm1 default) and `MQIMGREG` environment variable to be set. The value of `MQIMGREG` environment variable is the name of the docker container registry for MQ image.  

From the top level change to the `values` directory and run the build script.  
`cd values; export MQIMGREG=...; ./build.sh qm1`  

The build script creates `values/output` directory.  
```
ls values/output:
mqmodel.yaml
mqscic.yaml
qm1.env
qmini.yaml
qmspec.yaml
vault.yaml
webuser.yaml
```

`qm1.env` is queue manager environment file. It is used by the build script to create yaml files.

Other files are input to to the chart grouped by functionality.  

`qmspec.yaml` is required input for the chart. It configures core chart parameters.
`webuser.yaml` configures mq web console users. It is required if mq web console is enabled, otherwise optional.  
`qmini.yaml` is optional queue manager initialization file.
`qmscic.yaml` is opiontal mqsc command file applied to the queue manager on startup.
`mqmodel.yaml` is optional queue manager configuration file that is transformed to mqsc commands and applied at queue manager startup. Experimental.
`vault.yaml` is optional Hashi Corp vault integration configuration. Experimental.

**Values files customization**
Customize values files in the `values/output` directory.

Set ingress hostname parameters in the `qmspec.yaml` file.  
```
qmspec:
  ingress:
    qmgrHostname: "queue-manager-ingress-hostname"
    webcHostname: "web-console-ingress-hostname"
```

Note that there must be one to one correspondence between hostnames and chart release names.

**Required secrets and config maps**

MQ image pool secret:  
`kubectl create secret docker-registry image-pull-secret --docker-username=<u> --docker-password=<p> --docker-email=<e>`  

TLS secret:  
`kubectl create secret generic qm-tls --from-file=tls.key=</path/to/tls.key> --from-file=tls.crt=</path/to/tls.crt> --from-file=ca.crt=</path/to/ca.crt>`

Trust config map:  
`kubeclt create configmap qm-trust --from-file=ca1.crt=</path/to/ca1.crt> --from-file=ca2.crt=</path/to/ca2.crt> ...`

Ldap credentials secret:  
`kubectl create secret generic ldapcreds --from-literal=password=<ldappassword>`

**Chart installation**

Install mq helm chart with the `helm install` command from the top level git directory:  

`helm install -f values/output/qmspec.yaml [-f values/output/webuser.yaml] [-f values/output/mqscic.yaml] [-f values/output/qmini.yaml] [-f values/output/mqmodel.yaml] <release-name> mqchart/output`  

**Mqsc git integration**

**Examples and Reference**

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
