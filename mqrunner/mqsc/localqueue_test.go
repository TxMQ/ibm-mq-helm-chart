package mqsc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestLocalqueue_Mqsc(t *testing.T) {

	lq := Localqueue{
		Name:               "q.a",
		Descr:              "",
		Like:               "",
		Put:                "enabled",
		Get:                "enabled",
		DefaultPriority:    0,
		DefaultPersistence: false,
		Maxdepth:           0,
		Maxfsize:           0,
		Maxmsgl:            0,
		MsgDeliverySeq:     "priority", // priority|fifo
		acctq:              "",
		monq:               "",
		statq:              "",
		usage:              "",
		share:              false,
		Authority:          nil,
		Alter:              nil,
	}

	d, err := yaml.Marshal(&lq)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))

	mqsc := lq.Mqsc()
	fmt.Printf("%s\n", mqsc)

	data := `
name: q.b
descr: "application local queue"
put: true
get: true
defaultpriority: 2
defaultpersistence: false
maxdepth: 0
maxfsize: 0
maxmsgl: 0
msgdeliveryseq: priority
authority: []
alter: []
`
	lq = Localqueue{}
	err = yaml.Unmarshal([]byte(data), &lq)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", lq)

	mqsc = lq.Mqsc()
	fmt.Printf("%s\n", mqsc)
}
