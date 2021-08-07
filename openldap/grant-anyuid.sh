#!/bin/bash -x

oc create sa openldap

oc adm policy add-scc-to-user anyuid -z openldap --as system:admin
