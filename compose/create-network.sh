#!/bin/bash -x

net=qmnet

sudo podman network rm $net
sudo podman network create $net --driver bridge
