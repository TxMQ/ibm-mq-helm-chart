qmenv=$1

if [[ -z $qmenv ]]; then
echo qm environment file required: ./mqexplorer-mqsc-template.sh 'envfile'
exit 1
fi

# load environment
. $qmenv

cat <<EOF > output/mqexplorer-mqsc.yaml
mqscic: |
  *
  * there are 2 groups: admins, and readadmins who will use explorer channel
  * admins can perform any action, readers can view any object
  *

  * explorer channel
  DEFINE CHANNEL(MQ.EXPLORER) CHLTYPE(SVRCONN) SSLCIPH(ANY_TLS13) SSLCAUTH(OPTIONAL) REPLACE

  * channel rules
  SET CHLAUTH(MQ.EXPLORER) TYPE(ADDRESSMAP) USERSRC(CHANNEL) ADDRESS(*) ACTION(REPLACE)

  * connect to queue manager
  SET AUTHREC OBJTYPE(QMGR) GROUP(${ADMIN_GROUP}) AUTHADD(CONNECT, INQ, DSP)
  SET AUTHREC OBJTYPE(QMGR) GROUP(${READ_ADMIN_GROUP}) AUTHADD(CONNECT, INQ, DSP)

  SET AUTHREC PROFILE(SYSTEM.ADMIN.COMMAND.QUEUE) OBJTYPE(QUEUE) GROUP(${ADMIN_GROUP}) AUTHADD(PUT, INQ, DSP)
  SET AUTHREC PROFILE(SYSTEM.ADMIN.COMMAND.QUEUE) OBJTYPE(QUEUE) GROUP(${READ_ADMIN_GROUP}) AUTHADD(PUT, INQ, DSP)

  SET AUTHREC PROFILE(SYSTEM.MQEXPLORER.REPLY.MODEL) OBJTYPE(QUEUE) GROUP(${ADMIN_GROUP}) AUTHADD(GET, INQ, DSP)
  SET AUTHREC PROFILE(SYSTEM.MQEXPLORER.REPLY.MODEL) OBJTYPE(QUEUE) GROUP(${READ_ADMIN_GROUP}) AUTHADD(GET, INQ, DSP)

  * authreaders role
  SET AUTHREC PROFILE('**') OBJTYPE(QMGR) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(QUEUE) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(TOPIC) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(CHANNEL) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(PROCESS) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(NAMELIST) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(AUTHINFO) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(CLNTCONN) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(LISTENER) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(SERVICE) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)
  SET AUTHREC PROFILE('**') OBJTYPE(COMMINFO) GROUP(${READ_ADMIN_GROUP}) AUTHADD(DSP)

  * admins role
  SET AUTHREC PROFILE('**') OBJTYPE(QMGR) GROUP(${ADMIN_GROUP}) AUTHADD(ALL)
  SET AUTHREC PROFILE('**') OBJTYPE(QUEUE) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(TOPIC) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(CHANNEL) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(PROCESS) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(NAMELIST) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(AUTHINFO) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(CLNTCONN) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(LISTENER) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(SERVICE) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
  SET AUTHREC PROFILE('**') OBJTYPE(COMMINFO) GROUP(${ADMIN_GROUP}) AUTHADD(ALL, CRT)
EOF
