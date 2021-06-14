#!/bin/bash -x

# create mq directories
/opt/mqm/bin/crtmqdir -f -a

# create queue manager
/opt/mqm/bin/crtmqm -c "qm" -p 1414 -u SYSTEM.DEAD.LETTER.QUEUE -q qm

# start queue manager
/opt/mqm/bin/strmqm qm

qmpid=$!

# wait for queue manager
wait $qmpid
