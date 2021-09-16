package webmq

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"szesto.com/mqrunner/mqmodel"
	"szesto.com/mqrunner/util"
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

type Groupdef struct {
	ObjectClass string
	GroupNameAttr string
	GroupMembershipAttr string
}

type Userdef struct {
	ObjectClass string
	UsernameAttr string
}

type Connect struct {
	Realm string			// real name represents user registry
	Host string
	Port int
	Ldaptype string			// custom, mad, ...

	Binddn string
	Bindpassword string		// encode with securityUtility tool
	Basedn string			// start searching from this dn

	SslEnabled bool			// if enabled, ssl ref is required
}

type Ldapregistry struct {
	Connect Connect
	Groupdef Groupdef
	Userdef Userdef

	// example:
	// dn: cn=devs,ou=,o=,dc=,dc=
	// objectclass: groupOfNames
	// objectclass: top
	// cn=devs
	// member: uid=karson,ou=..
	// member: uid=roky,ou=...
	// member: uid=tobsky,ou=...
	//
	// groupfilter: (&(cn=%v)(objectclass=groupOfNames))
	// groupIdMap: *:cn - map a name of a group (devs) to an ldap entry by cn attribute type
	// groupMemberIdMap: groupOfNames:member - group membership attribute type
	// user filter: (&(uid=%v)(objectclass=inetOrgPerson))
	// userIdMap: *:uid - map a name of a user (roky) to an ldap entry by uid attribute type
	//
	// we can reuse mq ldap configuration attributes instead of filters and ldap entry maps
	// groupobjectclass, groupnameattr, groupmembershipattr
	// userobjectclass, usernameattr
	//
	// default custom filter:
	// (&(cn=%v)(|(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)
	//	(objectclass=groupOfURLs ???)))
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

func (webuser Webuser) ldapregistry() string {

	ldapf :=
	"   <ldapRegistry id=\"%s\"\n" +
	"      realm=\"%s\"\n" +
	"      host=\"%s\"\n" +
	"      port=\"%d\"\n" +
	"      bindDn=\"%s\"\n" +
	"      bindPassword=\"%s\"\n" +
	"      baseDn=\"%s\"\n" +
	"      ldapType=\"%s\"\n" +
	"      customFiltersRef=\"%s\">\n" +
	"   </ldapRegistry>\n\n"

	ldaptype := "custom"
	filtersid := "custom_filters"

	bindPassword := mqmodel.GetLdapBindPassword(webuser.Ldapregistry.Connect.Bindpassword)

	ldap := fmt.Sprintf(ldapf, "ldap", webuser.Ldapregistry.Connect.Realm, webuser.Ldapregistry.Connect.Host,
		webuser.Ldapregistry.Connect.Port, webuser.Ldapregistry.Connect.Binddn, bindPassword,
		webuser.Ldapregistry.Connect.Basedn, ldaptype, filtersid)

	filtersf :=
	"   <customLdapFilterProperties id=\"%s\"\n" +
	"      userFilter=\"%s\"\n" +
	"      groupFilter=\"%s\"\n" +
	"      userIdMap=\"%s\"\n" +
	"      groupIdMap=\"%s\"\n" +
	"      groupMemberIdMap=\"%s\">\n" +
	"   </customLdapFilterProperties>"

	userfilter := fmt.Sprintf("(&amp;(%s=%s)(objectclass=%s))", webuser.Ldapregistry.Userdef.UsernameAttr, "%v",
		webuser.Ldapregistry.Userdef.ObjectClass)

	groupfilter := fmt.Sprintf("(&amp;(%s=%s)(objectclass=%s))", webuser.Ldapregistry.Groupdef.GroupNameAttr, "%v",
		webuser.Ldapregistry.Groupdef.ObjectClass)

	useridmap := fmt.Sprintf("*:%s", webuser.Ldapregistry.Userdef.UsernameAttr)

	groupidmap := fmt.Sprintf("*:%s", webuser.Ldapregistry.Groupdef.GroupNameAttr)

	groupmemberidmap := fmt.Sprintf("%s:%s", webuser.Ldapregistry.Groupdef.ObjectClass,
		webuser.Ldapregistry.Groupdef.GroupMembershipAttr)

	filters := fmt.Sprintf(filtersf, filtersid, userfilter, groupfilter, useridmap, groupidmap, groupmemberidmap)

	return ldap + filters
}

func (webuser Webuser) tls(p12path, encpass string) string {

	ksf :=
	"   <keyStore id=\"webmqKeyStore\"\n"+
	"     location=\"%s\"\n" +
	"     type=\"PKCS12\"\n" +
	"     password=\"%s\"/>\n\n" +
	"   <ssl id=\"webmqSSLConfig\"\n" +
	"     clientAuthenticationSupported=\"true\"\n" +
	"     keyStoreRef=\"webmqKeyStore\"\n" +
	"     serverKeyAlias=\"default\"\n" +
	//"     trustStoreRef=\"defaultTrustStore\"\n" +
	"     sslProtocol=\"TLSv1.2\"/>\n\n" +
	"   <sslDefault sslRef=\"webmqSSLConfig\"/>"

	ks := fmt.Sprintf(ksf, p12path, encpass)
	return ks
}

func (webuser Webuser) Webuserxml(p12path, encpass string) string {

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
	ldap := webuser.ldapregistry()
	server += ldap + nl + nl

	// tls config
	if len(p12path) > 0 {
		tls := webuser.tls(p12path, encpass)
		server += tls + nl + nl
	}

	// client auth

	// variables
	for _, v := range webuser.Variables {
		server += webuser.variable(v.Name, v.Value) + nl
	}

	server += "</server>" + nl
	return server
}

func OutputWebuserxml(webconfigpath, webuserxmlpath, p12path, encpass string) error {

	data, err := ioutil.ReadFile(webconfigpath)
	if err != nil {
		return err
	}

	// unmarshall webuser config yaml
	webuser := Webuser{}

	err = yaml.Unmarshal([]byte(data), &webuser)
	if err != nil {
		return err
	}

	// p12 keystore path and encoded password are not part
	// of static webconfig and are passed in as parameters
	xml := webuser.Webuserxml(p12path, encpass)

	err = ioutil.WriteFile(webuserxmlpath, []byte(xml), 0777)
	if err != nil {
		return err
	}

	return nil
}

func ConfigureWebconsole() error {

	p12path := ""
	encpass := ""
	var err error = nil

	if util.IsEnableTls() {
		// import webconsole certs
		p12path, encpass, err = ImportWebconsoleCerts()
		if err != nil {
			return err
		}
	}

	// assume that web config is mounted at /etc/mqm/webuser/webuser.yaml
	// otherwise we can pass it as a parameter
	const webconfigpath = "/etc/mqm/webuser/webuser.yaml"

	// todo: use env var for directory path
	const webuserxmlpath = "/var/mqm/web/installations/Installation1/servers/mqweb/mqwebuser.xml"

	//if util.IsMultiInstance2() {
	//	logger.Logmsg(fmt.Sprintf("p-2, not merging web xml '%s' to '%s'", webconfigpath, webuserxmlpath))
	//	return nil
	//}

	err = OutputWebuserxml(webconfigpath, webuserxmlpath, p12path, encpass)
	if err != nil {
		return err
	}

	return nil
}
