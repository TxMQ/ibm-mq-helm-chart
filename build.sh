#!/bin/bash -x

. ./env.sh

registry=${1:-$MQIMGREG}

if [[ -z $registry ]]; then
echo mq image registry parameter required, either set MQIMGREG env or pass as argument: build.sh \<registry\>
exit 1
fi

# Chart.AppVersion value must match this value.
mqver=$MQVER

tag=$MQIMGTAG

image=txmq-mq-base-rpm-$mqver

sudo podman build --build-arg RPMDIR=$RPMDIR --build-arg MQVER=$mqver -t $image:$tag .

sudo podman tag "localhost/$image:$tag" "$registry/$image:$tag"

sudo podman push $registry/$image:$tag
