#!/bin/bash -xe

# copy mq runtime to install dir
noidir=$1
installdir=$2
uid=$3

./copy-mq-runtime.sh $noidir $installdir

# change installdir ownership
chown -R $uid:root $installdir/*

# accept license
$installdir/bin/mqlicense -accept

# create mout points
install -d -m 0755 -o $uid -g root /mnt
install -d -m 0755 -o $uid -g root /mnt/mqm
install -d -m 0755 -o $uid -g root /mnt/mqm/data
install -d -m 0755 -o $uid -g root /mnt/mqm-log
install -d -m 0755 -o $uid -g root /mnt/mqm-log/log
install -d -m 0755 -o $uid -g root /mnt/mqm-data
install -d -m 0755 -o $uid -g root /mnt/mqm-data/qmgrs

# mq configuration files
install -d -m 0775 -o $uid -g root /etc/mqm
install -d -m 0775 -o $uid -g root /etc/mqm/bin
install -d -m 0775 -o $uid -g root /etc/mqm/ssl
install -d -m 0775 -o $uid -g root /etc/mqm/pki/cert
install -d -m 0775 -o $uid -g root /etc/mqm/pki/trust

# directory for mq runtime files
install -d -m 0755 -o $uid -g root /run/runmqserver

touch /run/termination-log && chown $uid:root /run/termination-log

# link /var/mqm to /mnt/mqm/data
ln -s /mnt/mqm/data /var/mqm
