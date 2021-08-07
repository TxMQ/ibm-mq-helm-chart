package mqsc

import (
	"fmt"
	"strings"
)

type Authrec struct {
	Group []string
	Principal [] string
	Grant []string
	Revoke []string
}

func (authrec *Authrec) Mqsc(profile, objtype string) string {
	var authrecs []string

	// range over groups
	for _, group := range authrec.Group {

		// range over grants
		for _, grant := range authrec.Grant {
			t := "set authrec profile('%s') objtype(%s) group('%s') authadd(%s)"
			s := fmt.Sprintf(t, strings.ToUpper(profile), objtype, group, grant)
			authrecs = append(authrecs, s)
		}

		// range over revokes
		for _, revoke := range authrec.Revoke {
			t := "set authrec profile('%s') objtype(%s) group('%s') authrmv(%s)"
			s := fmt.Sprintf(t, strings.ToUpper(profile), objtype, group, revoke)
			authrecs = append(authrecs, s)
		}
	}

	// range over principals
	for _, principal := range authrec.Principal {

		// range over grants
		for _, grant := range authrec.Grant {
			t := "set authrec profile('%s') objtype(%s) principal('%s') authadd(%s)"
			s := fmt.Sprintf(t, strings.ToUpper(profile), objtype, principal, grant)
			authrecs = append(authrecs, s)
		}

		// range over revokes
		for _, revoke := range authrec.Revoke {
			t := "set authrec profile('%s') objtype(%s) principal('%s') authrmv(%s)"
			s := fmt.Sprintf(t, strings.ToUpper(profile), objtype, principal, revoke)
			authrecs = append(authrecs, s)
		}
	}

	return strings.Join(authrecs, "\n")
}
