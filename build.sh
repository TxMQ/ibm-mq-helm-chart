#!/bin/bash -x

. ./env.sh

registry=${1:-$MQIMGREG}

if [[ -z $registry ]]; then
echo output image registry parameter required
exit 1
fi

# Chart.AppVersion value must match this value.
mqver=$MQVER

tag=$MQIMGTAG

image=txmq-mq-base-rpm-$mqver

sudo podman build --build-arg RPMDIR=$RPMDIR --build-arg MQVER=$mqver -t $image:$tag .

sudo podman tag "localhost/$image:$tag" "$registry/$image:$tag"

sudo podman push $registry/$image:$tag
