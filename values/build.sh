#!/bin/bash -x

mkdir -p output

qmname=${1:-"qm1"}

./values-template.sh output $qmname

cp mqscic.yaml output
cp mq.yaml  output  
cp qmini.yaml output  
cp vault.yaml output
cp webuser.yaml output
