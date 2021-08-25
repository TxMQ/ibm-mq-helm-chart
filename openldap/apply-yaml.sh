#!/bin/bash -x

kubectl apply -f output/ldap-config-map.yaml
kubectl apply -f ./bitnami-sa.yaml
kubectl apply -f ./bitnami-service.yaml
kubectl apply -f ./bitnami-deployment.yaml
