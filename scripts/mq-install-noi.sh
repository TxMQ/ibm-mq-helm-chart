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