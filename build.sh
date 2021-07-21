#!/bin/bash -x

registry=$1

# derive version from the rpm directory
mqver="9.2.2.0"

tag="100"

image=txmq-mq-base-rpm-$mqver

sudo podman build --build-arg RPMDIR="rpm/MQServer" --build-arg MQVER=$mqver -t $image:$tag .

sudo podman tag "localhost/$image:$tag" "$registry/$image:$tag"

sudo podman push $registry/$image:$tag
