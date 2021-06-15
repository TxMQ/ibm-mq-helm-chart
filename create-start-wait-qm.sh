#!/bin/bash

set -x

/opt/mqm/bin/crtmqdir -f -a

/opt/mqm/bin/crtmqm -c "qm" -p 1414 -q qm

/opt/mqm/bin/strmqm qm

sleep 100d
