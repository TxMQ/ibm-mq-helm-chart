package mqsc

import "strings"

type Qmgr struct {
	Name string
	Access Chlauth
	Authority []Authrec
	Alter []string
}

func (qmgr *Qmgr) Mqsc() string {
	var qm []string

	// alter qmgr

	// chlauth
	s := qmgr.Access.Mqsc(qmgr.Name)
	qm = append(qm, s)

	for _, authrec := range qmgr.Authority {
		a := authrec.Mqsc(qmgr.Name, "qmgr")
		qm = append(qm, a)
	}

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