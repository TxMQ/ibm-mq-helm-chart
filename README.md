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

## Values.yaml file

There are a number of parameters that are required in values.yaml file.
Values.yaml file defines `.qmspec` and `.qmspec.qmconf` objects to configure queue manager.
`.webuser` object configures mq web console service.

`
qmspec:
  license:
    accept: true # required to accept license

  capabilities: mqbase # mq image capabilities, mqbase

  imagePullSecret: # docker registry secret to pull queue manager image
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

Queue manager running in kubernetes cluster must be configured to
authenticate against ldap server.

If you want to enable mqwebconsole it must be configured to
authenticate against the same ldap server as queue manager.

