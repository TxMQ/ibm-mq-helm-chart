#!/bin/bash -x

qmname=${1:-"qm1"}

if [[ -z $MQIMGREG ]]; then
echo mq image registry value required, set MQIMGREG env var
exit 1
fi

mkdir -p output
envfile=output/$qmname.env

if [[ ! -f $envfile ]]; then
./setenv.sh $qmname
fi

if [[ ! -f output/qmspec.yaml ]]; then
./qmspec-template.sh $envfile
fi

if [[ ! -f output/mqscic.yaml ]]; then
./mqscic-template.sh $envfile
fi

if [[ ! -f output/mqmodel.yaml ]]; then
./mqmodel-template.sh $envfile
fi

if [[ ! -f output/qmini.yaml ]]; then
cp qmini.yaml output  
fi

if [[ ! -f output/vault.yaml ]]; then
cp vault.yaml output
fi

if [[ ! -f output/webuser.yaml ]]; then
./webuser-template.sh $envfile
fi
