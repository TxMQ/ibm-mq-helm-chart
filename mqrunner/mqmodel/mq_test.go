package mqmodel

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMq_Mqsc(t *testing.T) {

	mq := Mq{
		Qmgr:       Qmgr{
			Name:   "",
			//Access: mqsc.Chlauth{},
			//Authority: []mqsc.Authrec{
			//	{
			//		Group:     nil,
			//		Principal: nil,
			//		Grant:     nil,
			//		Revoke:    nil,
			//	},
			//},
			Alter:     nil,
		},
		Auth: Auth{},
		//Svrconn:    []mqsc.Svrconn{
		//	{
		//		SvrconnProperties: mqsc.SvrconnProperties{
		//			Name:     "",
		//			Descr:    "",
		//			Maxmsgl:  0,
		//			Tls:      mqsc.Tls{
		//				Enabled:    false,
		//				ClientAuth: false,
		//				Ciphers:    nil,
		//			},
		//			discint:  0,
		//			hbint:    0,
		//			maxinst:  0,
		//			maxinstc: 0,
		//			sharecnv: 0,
		//		},
		//		Access: mqsc.Chlauth{},
		//		Authority:         []mqsc.Authrec{{
		//			Group:     nil,
		//			Principal: nil,
		//			Grant:     nil,
		//			Revoke:    nil,
		//		}},
		//		Alter: nil,
		//		allip: mqsc.Chlauth{},
		//	},
		//},
		//Localqueue: []mqsc.Localqueue{
		//	{
		//		Name:               "",
		//		Descr:              "",
		//		Like:               "",
		//		Put:                "enabled",
		//		Get:                "enabled",
		//		DefaultPriority:    0,
		//		DefaultPersistence: false,
		//		Maxdepth:           0,
		//		Maxfsize:           0,
		//		Maxmsgl:            0,
		//		MsgDeliverySeq:     "",
		//		Qtrigger:           mqsc.Qtrigger{},
		//		Qevents:            mqsc.Qevents{},
		//		Qcluster:           mqsc.Qcluster{},
		//		acctq:              "",
		//		monq:               "",
		//		statq:              "",
		//		usage:              "",
		//		share:              false,
		//		Authority:          nil,
		//		Alter:              nil,
		//	},
		//},
		Alter:      nil,
	}

	d, err := yaml.Marshal(&mq)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}

func TestMq_Mqsc2(t *testing.T) {
	data := `
qmgr:
  name: "qm1"
  alter: []

auth:
  ldap:
    connect:
      ldaphost: "openldap.mqmq.svc.cluster.local"
      ldapport: 389
      binddn: "cn=admin,dc=szesto,dc=com"
      bindpasswordsecret: "admin"
      tls: false
    groups:
      groupsearchbasedn: "ou=groups,dc=szesto.dc=com"
      objectclass: "groupOfNames"
      groupnameattr: "cn"
      groupmembershipattr: "member"
    users:
      usersearchbasedn: "ou=users,dc=szesto,dc=com"
      objectclass: "inetOrgPerson"
      usernameattr: "uid"
      shortusernameattr: "cn"

alter: []
`
	mq := Mq{}
	err := yaml.Unmarshal([]byte(data), &mq)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", mq)

	mqsc := mq.Mqsc()
	fmt.Printf("%v\n", mqsc)

}

func TestOutputmqsc(t *testing.T) {
	data := `
qmgr:
  name: "qm1"
  alter: []

auth:
  ldap:
    connect:
      ldaphost: "openldap.mqmq.svc.cluster.local"
      ldapport: 389
      binddn: "cn=admin,dc=szesto,dc=com"
      bindpassword: "admin"
      tls: false
    groups:
      groupsearchbasedn: "ou=groups,dc=szesto.dc=com"
      objectclass: "groupOfNames"
      groupnameattr: "cn"
      groupmembershipattr: "member"
    users:
      usersearchbasedn: "ou=users,dc=szesto,dc=com"
      objectclass: "inetOrgPerson"
      usernameattr: "uid"
      shortusernameattr: "cn"

alter: []
`
	// temp dir
	dir := os.TempDir()

	fmt.Printf("temp dir: %s\n", dir)

	// mqsc file name
	configyaml := filepath.Join(dir, "mqsctest.yaml")

	err := os.WriteFile(configyaml, []byte(data), 0777)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	mqscfile := filepath.Join(dir, "startup.mqsc")
	err = Outputmqsc(configyaml, mqscfile)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}
