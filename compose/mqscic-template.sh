qmenv=$1

if [[ -z $qmenv ]]; then
echo qm environment file required: ./mqscic-template.sh \<envfile\>
exit 1
fi

# load environment
. $qmenv

if [[ $LDAP_TYPE == "activedirectory" ]]; then
userobjectclass="USER"
usernameattr="sAMAccountName"
shortuser="employeeID"
groupobjectclass="GROUP"
groupnameattr="sAMAccountName"
groupmembershipattr="member"
else
# default: openldap
userobjectclass="inetOrgPerson"
usernameattr="uid"
shortuser="cn"
groupobjectclass="groupOfNames"
groupnameattr="cn"
groupmembershipattr="member"
fi

out=output/etc/mqm/mqsc

cat <<EOF > $out/qa.mqsc
DEFINE QLOCAL(Q.B) REPLACE DEFPSIST(YES)
EOF

cat <<EOF > $out/qb.mqsc
DEFINE QLOCAL(Q.C) REPLACE DEFPSIST(YES)
EOF

cat <<EOF > $out/authinfo.mqsc
ALTER QMGR CONNAUTH(USE.LDAP)

DEFINE AUTHINFO(USE.LDAP) + 
AUTHTYPE(IDPWLDAP) + 
ADOPTCTX(YES) + 
AUTHORMD(SEARCHGRP) + 
BASEDNG('$BASEDN_GROUPS') + 
BASEDNU('$BASEDN_USERS') + 
CLASSGRP($groupobjectclass) + 
CLASSUSR($userobjectclass) + 
CONNAME('$LDAP_HOST($LDAP_PORT)') + 
CHCKCLNT(required) + 
CHCKLOCL(optional) + 
DESCR('ldap authinfo') + 
FAILDLAY(1) + 
FINDGRP($groupmembershipattr) + 
GRPFIELD($groupnameattr) + 
LDAPPWD('admin') + 
LDAPUSER('$LDAP_USER') + 
NESTGRP(yes) + 
SECCOMM(no) + 
SHORTUSR($shortuser) + 
USRFIELD($usernameattr) + 
REPLACE

REFRESH SECURITY TYPE(CONNAUTH)
EOF
