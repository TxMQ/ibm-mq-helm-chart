package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"szesto.com/mqrunner/webmq"
)

func main() {

	// webconfig -ldap=custom|activedir -out filename
	
	// output webconfig template
	webuser := webmq.Webuser{
		Webroles: []webmq.Approle {
			{
				Name:   "MQWebAdmin",
				Users:  []string{},
				Groups: []string{},
			},
			{
				Name:   "MQWebAdminRO",
				Users:  []string{},
				Groups: []string{},
			},
			{
				Name:   "MQWebUser",
				Users:  []string{},
				Groups: []string{},
			},
		},
		Apiroles:     []webmq.Approle {
			{
				Name:   "MQWebAdmin",
				Users:  []string{},
				Groups: []string{},
			},
			{
				Name:   "MQWebAdminRO",
				Users:  []string{},
				Groups: []string{},
			},
			{
				Name:   "MQWebUser",
				Users:  []string{},
				Groups: []string{},
			},
		},
		Ldapregistry: webmq.Ldapregistry{
			Connect: webmq.Connect{
				Realm:        "realm",
				Host:         "host",
				Port:         389,
				Ldaptype:     "Custom",
				Binddn:       "cn=admin,dc=acme,dc=com",
				Bindpassword: "admin",
				Basedn:       "o=dev,dc=acme,dc=com",
				SslEnabled:   false,
			},
			Groupdef:     webmq.Groupdef{
				ObjectClass:         "groupOfNames",
				GroupNameAttr:       "cn",
				GroupMembershipAttr: "member",
			},
			Userdef:      webmq.Userdef{
				ObjectClass:  "inetOrgPerson",
				UsernameAttr: "uid",
			},
		},
		AllowedHosts: nil,
		Clientauth:   webmq.Clientauth{},
		Variables:    []webmq.Variable {
			{
				Name:  "httpsPort",
				Value: "9443",
			},
			{
				Name:  "httpHost",
				Value: "*",
			},
			{
				Name:  "mqRestCorsAllowedOrigints",
				Value: "*",
			},
		},
	}

	d, err := yaml.Marshal(&webuser)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
