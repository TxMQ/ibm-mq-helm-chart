#!/bin/bash -x

# install docker-compose

# start podman socket
sudo systemctl start podman.socket

# we should get OK
sudo curl -H "Content-Type: application/json" --unix-socket /var/run/docker.sock http://localhost/_ping
