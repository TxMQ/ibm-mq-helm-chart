#!/bin/bash -x

/opt/mqm/bin/runmqakm -cert -list -db key.kdb -stashed

/opt/mqm/bin/runmqakm -cert -details -db key.kdb -stashed -label ibmwebspheremqkarson

/opt/mqm/bin/runmqakm -cert -details -db key.kdb -stashed -label ca
