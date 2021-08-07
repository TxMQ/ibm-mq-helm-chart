apiVersion: build.openshift.io/v1
kind: BuildConfig
metadata:
  name: mq-server-noi
  namespace: mqmq
spec:
  source:
    git:
      uri: 'https://github.com/szesto/mq-operator.git'
    type: Git
    images:
    - from:
        kind: DockerImage
        name: "image-registry.openshift-image-registry.svc:5000/mqmq/mq-noi-tar:0.1"
      paths:
      - destinationDir: .
        sourcePath: /tmp/MQServer-noi
  strategy:
    dockerStrategy:
      dockerfilePath: Dockerfile-noi
  output:
    to:
      kind: DockerImage
      name: "image-registry.openshift-image-registry.svc:5000/mqmq/mq-server-noi:0.1"
