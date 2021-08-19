#!/bin/bash -x

qmname=${1:-"qm1"}

registry=${2:-$MQIMGREG}

if [[ -z $registry ]]; then
echo mq image registry parameter required, set MQIMGREG env, or pass on the command line: build.sh \<qm\> \<registry\>
exit 1
fi

mkdir -p output

if [[ ! -f output/values.yaml ]]; then
./values-template.sh output $qmname
fi

if [[ ! -f output/mqscic.yaml ]]; then
cp mqscic.yaml output
fi

if [[ ! -f output/mq.yaml ]]; then
cp mq.yaml  output  
fi

if [[ ! -f output/qmini.yaml ]]; then
cp qmini.yaml output  
fi

if [[ ! -f output/vault.yaml ]]; then
cp vault.yaml output
fi

if [[ ! -f output/webuser.yaml ]]; then
cp webuser.yaml output
fi
