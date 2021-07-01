#!/bin/bash -x

# directory owner
uid=$1

# create mout points
install -d -m 0755 -o $uid -g root /mnt
install -d -m 0755 -o $uid -g root /mnt/mqm
install -d -m 0755 -o $uid -g root /mnt/mqm/data
install -d -m 0755 -o $uid -g root /mnt/mqm-log
install -d -m 0755 -o $uid -g root /mnt/mqm-log/log
install -d -m 0755 -o $uid -g root /mnt/mqm-data
install -d -m 0755 -o $uid -g root /mnt/mqm-data/qmgrs

# mq init and config directories
install -d -m 0775 -o $uid -g root /etc/mqm
install -d -m 0777 -o $uid -g root /etc/mqm/sockets
install -d -m 0775 -o $uid -g root /etc/mqm/qmini
install -d -m 0775 -o $uid -g root /etc/mqm/mqsc
install -d -m 0775 -o $uid -g root /etc/mqm/bin
install -d -m 0775 -o $uid -g root /etc/mqm/ssl
install -d -m 0775 -o $uid -g root /etc/mqm/pki/cert
install -d -m 0775 -o $uid -g root /etc/mqm/pki/trust

# directory for mq runtime files
install -d -m 0755 -o $uid -g root /run/runmqserver

# termination log
touch /run/termination-log && chown $uid:root /run/termination-log

# link
ln -s /mnt/mqm/data /var/mqm
