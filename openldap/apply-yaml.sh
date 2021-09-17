#!/bin/bash

if [[ ! -f output/ldap-config-map.yaml ]]; then
echo Run \'./ldif-template.sh\' script to create \'output/ldif-config-map.yaml\' file.
exit 1
fi

set -x

kubectl apply -f output/ldif-config-map.yaml
kubectl apply -f ./bitnami-sa.yaml
kubectl apply -f ./bitnami-service.yaml
kubectl apply -f ./bitnami-deployment.yaml
