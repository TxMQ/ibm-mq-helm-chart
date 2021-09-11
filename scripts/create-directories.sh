#!/bin/bash -x

# directory owner
uid=$1

# create mout points
install -d -m 0755 -o $uid -g root /mnt
install -d -m 0755 -o $uid -g root /mnt/data
install -d -m 0755 -o $uid -g root /mnt/data/mqm
install -d -m 0755 -o $uid -g root /mnt/data/md
install -d -m 0755 -o $uid -g root /mnt/data/ld

# mq init and config directories
install -d -m 0775 -o $uid -g root /etc/mqm
install -d -m 0777 -o $uid -g root /etc/mqm/sockets
install -d -m 0775 -o $uid -g root /etc/mqm/qmini
install -d -m 0775 -o $uid -g root /etc/mqm/mqsc
install -d -m 0775 -o $uid -g root /etc/mqm/mqyaml
install -d -m 0775 -o $uid -g root /etc/mqm/bin
install -d -m 0775 -o $uid -g root /etc/mqm/ssl
install -d -m 0775 -o $uid -g root /etc/mqm/pki/cert
install -d -m 0775 -o $uid -g root /etc/mqm/pki/trust
install -d -m 0775 -o $uid -g root /etc/mqm/webuser

# directory for mq runtime files
install -d -m 0755 -o $uid -g root /run/runmqserver

# termination log
touch /run/termination-log && chown $uid:root /run/termination-log

# see if this works
cp /opt/mqm/bin/crtmqdir /opt/mqm/bin/crtmqdir_setuid && chown root:mqm /opt/mqm/bin/crtmqdir_setuid && chmod u+s /opt/mqm/bin/crtmqdir_setuid && chmod o+x /opt/mqm/bin/crtmqdir_setuid 

# link to /mnt/data/mqm from /var/mqm
rm -fr /var/mqm
ln -s /mnt/data/mqm /var
ln -s /mnt/data/md /var
ln -s /mnt/data/ld /var

# change link owner
chown $uid:root /var/mqm
chown $uid:root /var/md
chown $uid:root /var/ld
