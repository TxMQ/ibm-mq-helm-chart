#!/bin/bash -x

mkdir -p output

qmname=${1:-"qm1"}

./docker-compose-template.sh output $qmname

if [[ ! -d output/etc ]]; then
cp -r etc output
fi

if [[ ! -d output/ldif ]]; then
cp -r ldif output
fi
