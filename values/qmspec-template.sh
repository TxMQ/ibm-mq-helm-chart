#!/bin/bash

envfile=$1

if [[ -z $envfile ]]; then
echo qm env file required, ./qmspec-template.sh envfile
exit 1
fi

# load env
. $envfile

outdir=output

cat << EOF > $outdir/qmspec.yaml
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
  name: $QMNAME

  # custom docker image - required
  image: $MQIMGREG/txmq-mq-base-rpm-$MQVER:$MQIMGTAG

  # image pull policy IfNotPresent|Always
  imagePoolPolicy: Always

  service:

    # service type: ClusterIP, NodePort, LoadBalancer
    # if service type is ClusterIP, set ingress.enable=true
    # if service type is LoadBalancer, it is best to reserve ip address for the load balancer
    serviceType: ClusterIP

    qmgrPort: 1414
    httpswebcPort: 9443

  ingress:
    enable: "true"
    # nginx, openshift
    ingressClass: "nginx"

    # qmgr host name for the ingress.
    # typically qmgr name and a qualifier will be prepended to the ingress domain name.
    # qmgr host name must be unique for the chart release.
    # qmgr host name must be resolvable in naming service.
    # Ingress is configured in ssl passthrough mode. For this to work qmgr client code
    # must support SNI TLS extension and include qmgr hostname in sni extension.
    # if sni extension is not available, you can configure service type of LoadBalancer
    # and bypass ingress.
    qmgrHostname: ""

    # web console hostname for the ingress.
    # typically web console and a qualifier will be prepended to the ingress domain name.
    # web console host name must be unique for the chart release.
    # web console host name must be resolvable in naming service.
    # Ingress is configured in ssl passthrough mode. For this to work web browser
    # must support SNI TLS extension and include web console hostname in sni extension.
    # If sni extension is not available, you can configure service type of LoadBalancer
    # and bypass ingress.
    webcHostname: ""

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
    deleteClaim: false
    size: 2Gi
EOF
