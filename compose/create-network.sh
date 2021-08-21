#!/bin/bash

net=qmnet

# podman network exists return 0 is network exists

sudo podman network exists $net
if [[ $? == "1" ]]; then
   echo network $net does not exist

   echo creating network $net
   sudo podman network create $net --driver bridge
else
   echo network $net already exists
fi

