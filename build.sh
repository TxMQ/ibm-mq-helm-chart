#!/bin/bash -x

tag="0.17"

sudo podman build -t mq-adv-rpm .

sudo podman tag "localhost/mq-adv-rpm:latest" "docker.io/simong5000/txmq-mq-adv-rpm:$tag"

sudo podman push docker.io/simong5000/txmq-mq-adv-rpm:$tag

