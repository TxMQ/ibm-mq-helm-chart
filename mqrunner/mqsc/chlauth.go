package mqsc

import (
	"fmt"
	"strings"
)

type Chlauth struct {
	Defaultuser string
	Blockip []string
	Blockuser []string
	Allowip []string
}

func (chlauth *Chlauth) Mqsc(chlname string) string {
	var auths []string

	mcauser := chlauth.Defaultuser

	// range over blocked ips
	for _, ip := range chlauth.Blockip {
		t := "SET CHLAUTH(('%s') TYPE(addressmap) ADDRESS('%s') USERSRC(noaccess)"
		s := fmt.Sprintf(t, chlname, ip)
		auths = append(auths, s)
	}

	// range over blocked users
	for _, user := range chlauth.Blockuser {
		t := "SET CHLAUTH('%s') TYPE(blockuser) USERLIST('%s')"
		s := fmt.Sprintf(t, chlname, user)
		auths = append(auths, s)
	}

	// range allowed ips
	for _, ip := range chlauth.Allowip {

		s := ""
		if len(mcauser) > 0 {
			t := "set chlauth('%s') type(addressmap) address('%s') mcauser('%s') usersrc(map)"
			s = fmt.Sprintf(t, chlname, ip, mcauser)

		} else {
			t := "set chlauth('%s') type(addressmap) address('%s') usersrc(channel)"
			s = fmt.Sprintf(t, chlname, ip)
		}

		auths = append(auths, s)
	}

	return strings.Join(auths, "\n")
}
