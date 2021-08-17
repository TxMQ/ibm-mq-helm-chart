#!/bin/bash -x

mkdir -p output

./values-template.sh output

cp mqscic.yaml output
cp mq.yaml  output  
cp qmini.yaml output  
cp vault.yaml output
cp webuser.yaml output
