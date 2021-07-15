package mqsc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

	comment := fmt.Sprintf("%s", "ldap auth")
	mqsc = append(mqsc, formatComment(comment))

	ldap := mq.Auth.Ldap.Mqsc()
	mqsc = append(mqsc, ldap)

	comment = fmt.Sprintf("queue manager %s", mq.Qmgr.Name)
	mqsc = append(mqsc, formatComment(comment))

	qmgr := mq.Qmgr.Mqsc()
	mqsc = append(mqsc, qmgr)

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

func Outputmqsc(configyaml, mqscfile string) error {

	// read mq config yaml file
	data, err := ioutil.ReadFile(configyaml)
	if err != nil {
		if os.IsNotExist(err) {
			// nothing to do, return
			return nil
		}
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

	err = Printmqsc(mqsc, mqscfile)
	if err != nil {
		return err
	}

	return nil
}

func Printmqsc(mqsc, mqscfile string) error {
	//fmt.Printf("%s\n", mqsc)

	err := ioutil.WriteFile(mqscfile, []byte(mqsc), 0777)
	if err != nil {
		return err
	}

	return nil
}
