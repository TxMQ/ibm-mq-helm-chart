package webmq

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestWebuser_Webuserxml(t *testing.T) {

	wu := Webuser{
		Webroles: []Approle{
			{
				Name:   "MqWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MqWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MqWebUser",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MFTWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MFTWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Apiroles: []Approle{
			{
				Name:   "MqWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MqWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MqWebUser",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MFTWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MFTWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Ldapregistry: Ldapregistry{
			Realm:        "realm",
			Host:         "openldap.mqmq.svc.cluster.local",
			Port:         389,
			Ldaptype:     "Custom",
			Binddn:       "cn=admin,dc=szesto,dc=com",
			Bindpassword: "admin",
			Basedn:       "dc=szesto,dc=com",
			SslEnabled:   false,
			Groupdef:       Groupdef{
				ObjectClass:         "groupOfNames",
				GroupNameAttr:       "cn",
				GroupMembershipAttr: "member",
			},
			Userdef:        Userdef{
				ObjectClass:  "inetOrgPerson",
				UsernameAttr: "uid",
			},
		},
		Variables: []Variable{
			{
				Name:  "httpsPort",
				Value: "9443",
			},
			{
				Name:  "httpHost",
				Value: "*",
			},
			{
				Name:  "mqRestCorsAllowedOrigings",
				Value: "*",
			},
		},
	}

	xml := wu.Webuserxml("key.p12", "hello")
	fmt.Printf("%s\n", xml)
}

func TestWebuser_Webuserxml2(t *testing.T) {
	wu := Webuser{
		Webroles: []Approle{
			{
				Name:   "MQWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MQWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MQWebUser",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Apiroles: []Approle{
			{
				Name:   "MQWebAdmin",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MQWebAdminRO",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
			{
				Name:   "MQWebUser",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Ldapregistry: Ldapregistry{
			Realm:        "realm",
			Host:         "openldap.mqmq.svc.cluster.local",
			Port:         389,
			Ldaptype:     "Custom",
			Binddn:       "cn=admin,dc=szesto,dc=com",
			Bindpassword: "admin",
			Basedn:       "dc=szesto,dc=com",
			SslEnabled:   false,
			Groupdef:       Groupdef{
				ObjectClass:         "groupOfNames",
				GroupNameAttr:       "cn",
				GroupMembershipAttr: "member",
			},
			Userdef:        Userdef{
				ObjectClass:  "inetOrgPerson",
				UsernameAttr: "uid",
			},
		},
		Variables: []Variable{
			{
				Name:  "httpsPort",
				Value: "9443",
			},
			{
				Name:  "httpHost",
				Value: "*",
			},
			{
				Name:  "mqRestCorsAllowedOrigins",
				Value: "*",
			},
		},
	}

	d, err := yaml.Marshal(&wu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}