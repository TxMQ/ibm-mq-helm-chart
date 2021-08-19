qmenv=$1

if [[ -z $qmenv ]]; then
echo qm environment file required: ./mqscic-template.sh \<envfile\>
exit 1
fi

# load environment
. $qmenv

cat <<EOF > output/mqscic.yaml
mqscic: |
  DEFINE QLOCAL(Q.B) REPLACE DEFPSIST(YES)
  DEFINE QLOCAL(Q.C) REPLACE DEFPSIST(YES)
  *
  *ALTER QMGR CONNAUTH(USE.LDAP)
  *
  *DEFINE AUTHINFO(USE.LDAP) + 
  *AUTHTYPE(IDPWLDAP) + 
  *ADOPTCTX(YES) + 
  *AUTHORMD(SEARCHGRP) + 
  *BASEDNG('$BASEDN_GROUPS') + 
  *BASEDNU('$BASEDN_USERS') + 
  *CLASSGRP(groupOfNames) + 
  *CLASSUSR(inetOrgPerson) + 
  *CONNAME('$LDAP_HOST($LDAP_PORT)') + 
  *CHCKCLNT(required) + 
  *CHCKLOCL(optional) + 
  *DESCR('ldap authinfo') + 
  *FAILDLAY(1) + 
  *FINDGRP(member) + 
  *GRPFIELD(cn) + 
  *LDAPPWD('admin') + 
  *LDAPUSER('$LDAP_USER') + 
  *NESTGRP(yes) + 
  *SECCOMM(no) + 
  *SHORTUSR(cn) + 
  *USRFIELD(uid) + 
  *REPLACE
  *
  *REFRESH SECURITY TYPE(CONNAUTH)
EOF
