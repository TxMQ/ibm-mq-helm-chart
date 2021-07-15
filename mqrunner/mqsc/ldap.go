package mqsc

import "fmt"

type Auth struct {
	Ldap LdapAuthinfo
}

type LdapAuthinfo struct {
	Connect LdapConnect
	Groups  LdapGroups
	Users   LdapUsers
}

type LdapConnect struct {
	LdapHost string
	LdapPort int
	BindDn string
	BindPasswordSecret string
	Tls bool
}

type LdapGroups struct {
	GroupSearchBaseDn string
	ObjectClass string
	GroupNameAttr string
	GroupMembershipAttr string
}

type LdapUsers struct {
	UserSearchBaseDn string
	ObjectClass string
	UserNameAttr string
	ShortUserNameAttr string
}

const endl = "\n"
const cont = " + \n"
const star = "*\n"

func (ldap *LdapAuthinfo) Mqsc() string {

	t :=
		"alter qmgr connauth(use.ldap)" + endl +
		star +
		"define authinfo(use.ldap)" + cont +
		"authtype(IDPWLDAP)" + cont +
		"adoptctx(yes)" + cont +
		"authormd(searchgrp)" + cont +
		"basedng('%s')" + cont + // groups.groupSearchBaseDn
		"basednu('%s')" + cont + // users.userSearchBaseDn
		"CLASSGRP(%s)" + cont + // {{ .groups.objectClass }})" +
		"CLASSUSR(%s)" + cont + // {{ .users.objectClass }}) +
		"CONNAME('%s(%d)')" + cont + // {{ .connect.ldapHost }}({{ .connect.ldapPort }})') +
		"CHCKCLNT(required)" + cont +
		"CHCKLOCL(optional)" + cont +
		"DESCR('ldap authinfo')" + cont +
		"FAILDLAY(1)" + cont +
		"FINDGRP(%s)" + cont + // ({{ .groups.groupMembershipAttr }}) +
		"GRPFIELD(%s)" + cont + // {{ .groups.groupNameAttr }}) +
		"LDAPPWD('%s')" + cont + // {{ .connect.bindPasswordSecret | squote }}) +
		"LDAPUSER('%s')" + cont + // {{ .connect.bindDn | squote }}) +
		"NESTGRP(yes)" + cont  +
		"SECCOMM(no)" + cont + // todo: parameterize, ssl to ldap
		"SHORTUSR(%s)" + cont + // {{ .users.shortUserNameAttr }}) +
		"USRFIELD(%s)" + cont + // {{ .users.userNameAttr }});
		"replace" + endl +
		star +
		"REFRESH SECURITY TYPE(CONNAUTH)" + endl

	mqsc := fmt.Sprintf(t, ldap.Groups.GroupSearchBaseDn, ldap.Users.UserSearchBaseDn,
		ldap.Groups.ObjectClass, ldap.Users.ObjectClass,
		ldap.Connect.LdapHost, ldap.Connect.LdapPort,
		ldap.Groups.GroupMembershipAttr, ldap.Groups.GroupNameAttr,
		ldap.Connect.BindPasswordSecret, ldap.Connect.BindDn,
		ldap.Users.ShortUserNameAttr, ldap.Users.UserNameAttr)

	return mqsc
}
