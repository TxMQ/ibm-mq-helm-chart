#!/bin/bash -x

registry=$1

# Chart.AppVersion value must match this value.
mqver="9.2.2.0"

tag="159"

image=txmq-mq-base-rpm-$mqver

sudo podman build --build-arg RPMDIR="rpm/MQServer" --build-arg MQVER=$mqver -t $image:$tag .

sudo podman tag "localhost/$image:$tag" "$registry/$image:$tag"

sudo podman push $registry/$image:$tag
