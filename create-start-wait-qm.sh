#!/bin/bash -x

# create mq directories
/opt/mqm/bin/crtmqdir -f -a

# create queue manager
/opt/mqm/bin/crtmqm -c "qm" -p 1414 -u SYSTEM.DEAD.LETTER.QUEUE -q qm

# start queue manager
/opt/mqm/bin/strmqm qm

pid=`ps -ef | grep qm | grep amqzxma0 | tr -s " " | cut -d " " -f 2`

ps -ef | grep $pid

sleep 10d

