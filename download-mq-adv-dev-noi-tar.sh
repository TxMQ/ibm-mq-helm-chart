#!/bin/bash -xe

tardir=$1

url="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqadv"
tarfile="9.2.2.0-IBM-MQ-Advanced-for-Developers-Non-Install-LinuxX64.tar.gz"

wget -q "$url/$tarfile"

tar xzvf $tarfile -C $tardir
