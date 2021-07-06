package mqsc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestSvrconn_Mqsc(t *testing.T) {

	data := `
svrconnproperties:
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
`
	// todo: authorities applicable to the channel

	svrconn := Svrconn{}

	err := yaml.Unmarshal([]byte(data), &svrconn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", svrconn)

	mqsc := svrconn.Mqsc()
	fmt.Printf("%s\n", mqsc)

	svrconn = Svrconn{
		SvrconnProperties: SvrconnProperties{
			Name:     "app.channel",
			Descr:    "app channel",
			Tls:      Tls{
				Enabled:    false,
				ClientAuth: false,
				Ciphers:    nil,
			},
			discint:  0,
			hbint:    0,
			maxinst:  0,
			maxinstc: 0,
			Maxmsgl:  0,
			sharecnv: 0,
		},
		Access:            Chlauth{},
		Authority:         nil,
		Alter:             nil,
		allip:             Chlauth{},
	}

	d, err := yaml.Marshal(&svrconn)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))

	mqsc = svrconn.Mqsc()
	fmt.Printf("%s\n", mqsc)
}
