#!/bin/bash -x

# prepare volumes for mq container

img="registry.access.redhat.com/ubi8/ubi-minimal:8.4"

cidfile=etcmqm.cid

sudo podman volume rm mqsc
sudo podman volume rm qmtls
sudo podman volume rm qmtrust
sudo podman volume rm webuser
sudo podman volume rm ldif

sudo podman run --cidfile $cidfile --name etcmqm -v mqsc:/etc/mqm/mqsc -v qmtls:/etc/mqm/pki/cert -v qmtrust:/etc/mqm/pki/trust -v webuser:/etc/mqm/webuser -v ldif:/ldifs $img /bin/sh

cid=$(cat $cidfile)

# mqsc volume
for f in `ls etc/mqm/mqsc/*.mqsc`
do 
sudo podman cp $f etcmqm:/etc/mqm/mqsc
done

# qmtls volume
for f in `ls etc/mqm/pki/cert/*`
do
sudo podman cp $f etcmqm:/etc/mqm/pki/cert
done

# qmtrust volume
for f in `ls etc/mqm/pki/trust/*`
do
sudo podman cp $f etcmqm:/etc/mqm/pki/trust
done

# webuser volume
for f in `ls etc/mqm/webuser/webuser.yaml`
do
sudo podman cp $f etcmqm:/etc/mqm/webuser
done

# openldap ldif volume
for f in `ls ldif/*.ldif`
do
sudo podman cp $f etcmqm:/ldifs
done

sudo podman container stop $cid
sudo podman rm $cid
rm -f $cidfile

# to find voldir: sudo podman inspect mqsc
voldir="/var/lib/containers/storage/volumes/mqsc/_data"
for f in `sudo ls $voldir`
do
sudo chown 1001:1001 $voldir/$f
#sudo cat $voldir/$f
done

voldir="/var/lib/containers/storage/volumes/qmtls/_data"
for f in `sudo ls $voldir`
do
sudo chown 1001:1001 $voldir/$f
#sudo cat $voldir/$f
done

voldir="/var/lib/containers/storage/volumes/qmtrust/_data"
for f in `sudo ls $voldir`
do
sudo chown 1001:1001 $voldir/$f
#sudo cat $voldir/$f
done

voldir="/var/lib/containers/storage/volumes/webuser/_data"
for f in `sudo ls $voldir`
do
sudo chown 1001:1001 $voldir/$f
#sudo cat $voldir/$f
done
