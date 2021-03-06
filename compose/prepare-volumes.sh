#!/bin/bash -x

# prepare volumes for mq container

img="registry.access.redhat.com/ubi8/ubi-minimal:8.4"

cidfile=etcmqm.cid

sudo podman volume rm mqsc
sudo podman volume rm mqyaml
sudo podman volume rm qmtls
sudo podman volume rm qmtrust
sudo podman volume rm webuser
sudo podman volume rm qmini
sudo podman volume rm ldif
sudo podman volume rm mqmq
sudo podman volume rm mqmd
sudo podman volume rm mqld

sudo podman run --cidfile $cidfile --name etcmqm -v mqsc:/etc/mqm/mqsc -v mqyaml:/etc/mqm/mqyaml -v qmtls:/etc/mqm/pki/cert -v qmtrust:/etc/mqm/pki/trust -v webuser:/etc/mqm/webuser -v qmini:/etc/mqm/qmini -v ldif:/ldifs $img /bin/sh

cid=$(cat $cidfile)

# mqsc volume
for f in `ls output/etc/mqm/mqsc/*.mqsc`
do 
sudo podman cp $f etcmqm:/etc/mqm/mqsc
done

# mqmodel volume
for f in `ls output/etc/mqm/mqyaml/*.yaml`
do 
sudo podman cp $f etcmqm:/etc/mqm/mqyaml
done

# qmtls volume
for f in `ls output/etc/mqm/pki/cert/*`
do
sudo podman cp $f etcmqm:/etc/mqm/pki/cert
done

# qmtrust volume
for f in `ls output/etc/mqm/pki/trust/*`
do
sudo podman cp $f etcmqm:/etc/mqm/pki/trust
done

# webuser volume
for f in `ls output/etc/mqm/webuser/webuser.yaml`
do
sudo podman cp $f etcmqm:/etc/mqm/webuser
done

# qmini volume
for f in `ls output/etc/mqm/qmini/*.yaml`
do 
sudo podman cp $f etcmqm:/etc/mqm/qmini
done

# openldap ldif volume
for f in `ls output/ldif/*.ldif`
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

voldir="/var/lib/containers/storage/volumes/mqyaml/_data"
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

voldir="/var/lib/containers/storage/volumes/qmini/_data"
for f in `sudo ls $voldir`
do
sudo chown 1001:1001 $voldir/$f
#sudo cat $voldir/$f
done
