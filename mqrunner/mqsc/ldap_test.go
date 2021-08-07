package mqsc

import (
	"fmt"
	"log"
	"testing"
)
import "gopkg.in/yaml.v2"

func TestLdapAuthinfo_Mqsc(t *testing.T) {

	data := `
connect:
  ldaphost: openldap.mqmq.svc.cluster.local
  ldapport: 389
  binddn: cn=admin,dc=szesto,dc=com
  bindpassword: admin
  tls: false
groups:
  groupsearchbasedn: ou=groups,dc=szesto,dc=com
  objectclass: groupOfNames
  groupnameattr: cn
  groupmembershipattr: member
users:
  usersearchbasedn: ou=users,dc=szesto,dc=com
  objectclass: inetOrgPerson
  usernameattr: uid
  shortusernameattr: cn
`
	var ldapa LdapAuthinfo

	err := yaml.Unmarshal([]byte(data), &ldapa)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", ldapa)

	mqsc := ldapa.Mqsc()
	fmt.Printf("%v\n", mqsc)

	ldapauth := LdapAuthinfo{
		Connect: LdapConnect{
			LdapHost:           "openldap.mqmq.svc.cluster.local",
			LdapPort:           389,
			BindDn:             "cn=admin,dc=szesto,dc=com",
			BindPassword: 		"admin",
			Tls:                false,
		},
		Groups:  LdapGroups{
			GroupSearchBaseDn:   "ou=groups,dc=szesto,dc=com",
			ObjectClass:         "groupOfNames",
			GroupNameAttr:       "cn",
			GroupMembershipAttr: "member",
		},
		Users:   LdapUsers{
			UserSearchBaseDn:  "ou=users,dc=szesto,dc=com",
			ObjectClass:       "inetOrgPerson",
			UserNameAttr:      "uid",
			ShortUserNameAttr: "cn",
		},
	}

	d, err := yaml.Marshal(&ldapauth)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
