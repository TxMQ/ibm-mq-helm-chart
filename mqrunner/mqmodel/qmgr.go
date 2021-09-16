package mqmodel

import (
	"strings"
)

type Qmgr struct {
	Name      string
	//Access    mqsc.Chlauth
	//Authority []mqsc.Authrec
	Alter     []string
}

func (qmgr *Qmgr) Mqsc() string {
	var qm []string

	// chlauth
	//s := qmgr.Access.Mqsc(qmgr.Name)
	//qm = append(qm, s)
	//
	//for _, authrec := range qmgr.Authority {
	//	a := authrec.Mqsc(qmgr.Name, "qmgr")
	//	qm = append(qm, a)
	//}

	// alter qmgr
	for _, alt := range qmgr.Alter {
		a := qmgr.alter(alt)
		qm = append(qm, a)
	}

	return strings.Join(qm, "\n")
}

func (qmgr *Qmgr) alter(alt string) string {

	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(alt)), "alter") {
		return alt
	} else {
		return "alter qmgr " + alt
	}
}