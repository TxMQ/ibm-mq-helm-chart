#!/bin/bash -x

. ./env.sh

registry=${1:-$MQIMGREG}

# Chart.AppVersion value must match this value.
mqver=$MQVER

tag=$MQIMGTAG

image=txmq-mq-base-rpm-$mqver

sudo podman build --build-arg RPMDIR=$RPMDIR --build-arg MQVER=$mqver -t $image:$tag .

sudo podman tag "localhost/$image:$tag" "$registry/$image:$tag"

sudo podman push $registry/$image:$tag
