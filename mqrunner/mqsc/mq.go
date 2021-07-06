package mqsc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Mq struct {
	Qmgr Qmgr
	Auth Auth
	Svrconn []Svrconn
	Localqueue []Localqueue
	Alter []string
}

func (mq *Mq) Mqsc() string {

	var mqsc []string

	comment := fmt.Sprintf("queue manager %s", mq.Qmgr.Name)
	mqsc = append(mqsc, formatComment(comment))

	qmgr := mq.Qmgr.Mqsc()
	mqsc = append(mqsc, qmgr)

	comment = fmt.Sprintf("%s", "ldap auth")
	mqsc = append(mqsc, formatComment(comment))

	ldap := mq.Auth.Ldap.Mqsc()
	mqsc = append(mqsc, ldap)

	for _, svrconn := range mq.Svrconn {
		comment = fmt.Sprintf("svrconn channel %s", svrconn.Name)
		mqsc = append(mqsc, formatComment(comment))

		svrc := svrconn.Mqsc()
		mqsc = append(mqsc, svrc)
	}

	for _, lq := range mq.Localqueue {
		comment = fmt.Sprintf("local queue %s", lq.Name)
		mqsc = append(mqsc, formatComment(comment))

		qs := lq.Mqsc()
		mqsc = append(mqsc, qs)
	}

	comment = fmt.Sprintf("%s", "alter statements")
	mqsc = append(mqsc, formatComment(comment))

	for _, alt := range mq.Alter {
		mqsc = append(mqsc, alt)
	}

	return strings.Join(mqsc, "\n")
}

func formatComment(comment string) string {
	return fmt.Sprintf("\n*\n* %s\n*\n", comment)
}

func Outputmqsc(mqscfile string) error {

	// read mq config yaml file
	data, err := ioutil.ReadFile(mqscfile)
	if err != nil {
		return err
	}

	// unmarshall config yaml
	mq := Mq{}

	err = yaml.Unmarshal(data, &mq)
	if err != nil {
		return err
	}

	// output mqsc
	mqsc := mq.Mqsc()
	Printmqsc(mqsc)

	return nil
}

func Printmqsc(mqsc string) {
	fmt.Printf("%s\n", mqsc)
}