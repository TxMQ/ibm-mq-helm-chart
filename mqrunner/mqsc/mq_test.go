package mqsc

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
			Name:      "",
			Access:    Chlauth{},
			Authority: []Authrec{
				{
					Group:     nil,
					Principal: nil,
					Grant:     nil,
					Revoke:    nil,
				},
			},
			Alter:     nil,
		},
		Auth:       Auth{},
		Svrconn:    []Svrconn {
			{
				SvrconnProperties: SvrconnProperties{
					Name:     "",
					Descr:    "",
					Maxmsgl:  0,
					Tls:      Tls{
						Enabled:    false,
						ClientAuth: false,
						Ciphers:    nil,
					},
					discint:  0,
					hbint:    0,
					maxinst:  0,
					maxinstc: 0,
					sharecnv: 0,
				},
				Access:            Chlauth{},
				Authority:         []Authrec{{
					Group:     nil,
					Principal: nil,
					Grant:     nil,
					Revoke:    nil,
				}},
				Alter:             nil,
				allip:             Chlauth{},
			},
		},
		Localqueue: []Localqueue{
			{
				Name:               "",
				Descr:              "",
				Like:               "",
				Put:                false,
				Get:                false,
				DefaultPriority:    0,
				DefaultPersistence: false,
				Maxdepth:           0,
				Maxfsize:           0,
				Maxmsgl:            0,
				MsgDeliverySeq:     "",
				Qtrigger:           Qtrigger{},
				Qevents:            Qevents{},
				Qcluster:           Qcluster{},
				acctq:              "",
				monq:               "",
				statq:              "",
				usage:              "",
				share:              false,
				Authority:          nil,
				Alter:              nil,
			},
		},
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
  access:
    defaultuser: ""
    blockip: []
    blockuser: []
    allowip: []
  authority:
  - group: [devs]
    principal: [karson]
    grant: [connect]
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

svrconn:
- svrconnproperties:
    name: app.channel
    maxmsgl: 4096
  tls:
    enabled: false
    clientauth: false
    ciphers: []
  access:
    defaultuser: ""
    blockip: ["10.2.100.*"]
    blockuser: []
    allowip: []
  authority: []
  alter:
  - monchl(off)
  - alter channel(app.channel) monchl(low)

localqueue:
- name: q.a
  put: true
  get: true
  defaultprioprity: 2
  defaultpersistence: true
- name: q.b
  like: q.a

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
  access:
    defaultuser: ""
    blockip: []
    blockuser: []
    allowip: []
  authority:
  - group: [devs]
    principal: [karson]
    grant: [connect]
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

svrconn:
- svrconnproperties:
    name: app.channel
    maxmsgl: 4096
  tls:
    enabled: false
    clientauth: false
    ciphers: []
  access:
    defaultuser: ""
    blockip: ["10.2.100.*"]
    blockuser: []
    allowip: []
  authority: []
  alter:
  - monchl(off)
  - alter channel(app.channel) monchl(low)

localqueue:
- name: q.a
  put: true
  get: true
  defaultprioprity: 2
  defaultpersistence: true
- name: q.b
  like: q.a

alter: []
`
	// temp dir
	dir := os.TempDir()

	fmt.Printf("temp dir: %s\n", dir)

	// mqsc file name
	mqscfile := filepath.Join(dir, "mqsctest.yaml")

	err := os.WriteFile(mqscfile, []byte(data), 0777)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	err = Outputmqsc(mqscfile)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}
