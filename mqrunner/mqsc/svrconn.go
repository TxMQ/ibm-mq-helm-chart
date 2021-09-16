package mqsc

import (
	"fmt"
	"strings"
	"szesto.com/mqrunner/mqmodel"
)

type Tls struct {
	Enabled bool
	ClientAuth bool
	Ciphers []string
}

type SvrconnProperties struct {
	Name string
	Descr string
	Maxmsgl int // // max message length that can be transmitted on the channel; limited by maxmsgl on the qm
	Tls Tls

	// use alter channel statement to change unexported fields
	discint int 	// min time in seconds for which channel remains active wo client communication
	hbint int 		// 0 - 999999; 300; interval in seconds between mca heartbeats
	maxinst int 	// 0 - 999999999; the max number of individual simultanious instances of svrconn
	maxinstc int 	// 0 - 999999999; the max number of simultanious individual channels from the same client
	sharecnv int 	// max number of conversations that can be shared on channel instance. 0-999999999
}

type Svrconn struct {
	SvrconnProperties

	// related mqsc
	Access Chlauth
	Authority []Authrec
	Alter []string
	allip Chlauth
}

func (props *SvrconnProperties) mqsc() string {

	// set unexported fields defaults
	props.discint = 0
	props.hbint = 300
	props.maxinst = 1_000_000
	props.maxinstc = 100_000
	props.sharecnv = 100_000

	// set exported defaults
	if len(props.Descr) == 0 {
		props.Descr = fmt.Sprintf("srvconn channel %s", strings.ToUpper(props.Name))
	}

	t :=	"define channel('%s')" + mqmodel.cont + // name
			"chltype(svrconn)" + mqmodel.cont +
			"descr('%s')" + mqmodel.cont + // descr
			"trptype(tcp)" + mqmodel.cont +
			"monchl(qmgr)" + mqmodel.cont +
			"discint(%d)" + mqmodel.cont + // discint
			"hbint(%d)" + mqmodel.cont + // hbint
			"maxinst(%d)" + mqmodel.cont + // maxinst
			"maxinstc(%d)" + mqmodel.cont + // maxinstc
			"maxmsgl(%d)" + mqmodel.cont + // maxmsgl
			"sharecnv(%d)" + mqmodel.endl // sharecnv

	s := fmt.Sprintf(t, strings.ToUpper(props.Name), props.Descr,
		props.discint, props.hbint, props.maxinst, props.maxinstc, props.Maxmsgl, props.sharecnv)

	return s
}

func (svrconn *Svrconn) Mqsc() string {

	var mqsc []string

	// channel properties
	p := svrconn.mqsc()
	mqsc = append(mqsc, p)

	// channel derived mqsc
	s := svrconn.Access.Mqsc(svrconn.Name)
	if len(s) > 0 {
		mqsc = append(mqsc, s)
	}

	for _, authrec := range svrconn.Authority {
		a := authrec.Mqsc(svrconn.Name, "channel")
		if len(a) > 0 {
			mqsc = append(mqsc, a)
		}
	}

	// todo: tls

	// allip
	svrconn.allip = Chlauth{
		Allowip: []string {"*"},
	}

	s = svrconn.allip.Mqsc(svrconn.Name)
	mqsc = append(mqsc, s)

	// alter channel
	for _, alt := range svrconn.Alter {
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(alt)), "alter") {
			mqsc = append(mqsc, alt)

		} else {
			a := fmt.Sprintf("alter channel(%s) %s", svrconn.Name, alt)
			mqsc = append(mqsc, a)
		}
	}

	return strings.Join(mqsc, "\n")
}
