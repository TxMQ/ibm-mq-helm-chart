package mqmodel

import (
	"fmt"
	"log"
	mqsc2 "szesto.com/mqrunner/mqsc"
	"testing"
)
import "gopkg.in/yaml.v2"

func TestQmgr_Mqsc(t *testing.T) {

	var data = `
name: qm1
access:
  defaultuser: nobody
  blockip: [10.5.*, 192.168.2.*]
  blockuser: [zorro]
authority:
- group: [devs]
  principal: [karson]
  grant: [connect]
alter:
- chlauth(enabled)
- alter qmgr chlauth(disabled)
`
	qmgr := Qmgr{}
	err := yaml.Unmarshal([]byte(data), &qmgr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", qmgr)

	mqsc := qmgr.Mqsc()
	fmt.Printf("%v\n", mqsc)

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m:\n%v\n\n", m)

	qm1 := Qmgr{
		Name:      "qm1",
		Access:  mqsc2.Chlauth{
			Defaultuser: "",
			Blockip:     []string{"10.5.*"},
			Blockuser:   []string{"zorro"},
		},
		Authority: []mqsc2.Authrec{{
			Group:  []string{"devs"},
			Principal: []string{"karson"},
			Grant:     []string{"connect"},
			Revoke:    nil,
		}},
	}

	d, err := yaml.Marshal(&qm1)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
