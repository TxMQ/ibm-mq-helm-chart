#!/bin/bash

sa=$1

if [[ -z $sa ]]; then
echo service account name required: grant-anyuid.sh \<serviceaccount\>
exit 1
fi

oc adm policy add-scc-to-user anyuid -z $sa --as system:admin
