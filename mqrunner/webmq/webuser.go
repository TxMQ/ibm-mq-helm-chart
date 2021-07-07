package webmq

import (
	"fmt"
	"strings"
)

type Approle struct {
	Name   string
	Users  []string
	Groups []string
}

type Variable struct {
	Name  string
	Value string
}

type Ldapregistry struct {
	Realm string
	Host string
	Port int
	Ldaptype string
	Binddn string
	Bindpassword string
	Basedn string
	Userfilter string
	Groupfilter string
}

type Clientauth struct {
	Keystorepass string
	Truststorepass string
	Enabled bool
}

type Webuser struct {
	Webroles []Approle
	Apiroles []Approle
	Ldapregistry Ldapregistry
	AllowedHosts []string // hosts allowed to connect to mqweb
	Clientauth Clientauth
	// replace variables with keys
	Variables []Variable
}

func (webuser Webuser) xmldecl() string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
}

func (webuser Webuser) features() string {

	f :=
	"   <featureManager>\n" +
	"      <feature>appSecurity-2.0</feature>\n" +
	"      <feature>ldapRegistry-3.0</feature>\n" +
	"      <feature>basicAuthenticationMQ-1.0</feature>\n" +
	"   </featureManager>"

	return f
}

func (webuser Webuser) variable(name, value string) string {
	return fmt.Sprintf("   <variable name=\"%s\" value=\"%s\"/>", name, value)
}

func (webuser Webuser) virtualhosts(aliases []string) string {
	f :=
	"   <virtualHost allowFromEndpointRef=\"defaultHttpEndpoint\" id=\"default_host\">\n" +
	"      %s\n" +
	"   </virtualHost>"

	var hlist []string
	for _, alias := range aliases {
		hlist = append(hlist, fmt.Sprintf("<hostAlias>%s</hostAlias>", alias))
	}

	vh := fmt.Sprintf(f, strings.Join(hlist, "\n"))
	return vh
}

func (webuser Webuser) mapGroupToRole(role, group, realm string) string {
	secrole :=
	"         <security-role name=\"%s\">\n" +
	"            <group name=\"%s\" realm=\"%s\"/>\n" +
	"         </security-role>"
	return fmt.Sprintf(secrole, role, group, realm)
}

func (webuser Webuser) mapUserToRole(role, user, realm string) string {
	secrole :=
	"         <security-role name=\"%s\">\n" +
	"            <user name=\"%s\" realm=\"%s\"/>\n" +
	"         </security-role>"
	return fmt.Sprintf(secrole, role, user, realm)
}

func (webuser Webuser) roleBindings(roles []Approle) string {

	appbinding :=
	"   <enterpriseApplication id=\"com.ibm.mq.console\">\n" +
	"      <application-bnd>\n"

	for _, role := range roles {
		// Groups
		for _, group := range  role.Groups {
			rolemap := webuser.mapGroupToRole(role.Name, group, "realm")
			appbinding += rolemap + "\n"
		}

		// Users
		for _, user := range role.Users {
			rolemap := webuser.mapUserToRole(role.Name, user, "realm")
			appbinding += rolemap + "\n"
		}
	}

	appbinding +=
	"      </application-bnd>\n" +
	"   </enterpriseApplication>"

	return appbinding
}

func (webuser Webuser) Webuserxml() string {

	const nl = "\n"

	server := webuser.xmldecl() + nl
	server += "<server>" + nl

	server += webuser.features() + nl + nl

	// roles for mq console
	consoleRoles := webuser.roleBindings(webuser.Webroles)
	server += consoleRoles + nl + nl

	// roles for mq rest api
	apiRoles := webuser.roleBindings(webuser.Apiroles)
	server += apiRoles + nl + nl

	// ldap registry

	// tls config

	// client auth

	// variables
	for _, v := range webuser.Variables {
		server += webuser.variable(v.Name, v.Value) + nl
	}

	server += "</server>" + nl
	return server
}