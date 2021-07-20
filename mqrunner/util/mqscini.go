package util

import (
	"io/ioutil"
	"log"
	"os"
	"szesto.com/mqrunner/mqsc"
)

const _mqyaml = "/etc/mqm/mqyaml/mq.yaml"
const _mqyamlout = "/etc/mqm/mqyamlout.mqsc"
const _mqscini = "/etc/mqm/mqsc/mqscini.mqsc"
const _mqscic = "/etc/mqm/mqscic.mqsc"
const _qmini = "/etc/mqm/qmini/qmini.ini"

func MergeMqscFiles() error {

	// transform mqconfig yaml into mqsc output file
	mqyaml := _mqyaml
	mqyamlout := _mqyamlout

	ismqyaml := true

	_, err := os.Stat(mqyaml)
	if err != nil && os.IsNotExist(err) {
		log.Printf("merge-mqsc-files: mq yaml file '%s' not found\n", mqyaml)
		ismqyaml = false

	} else if err != nil {
		return err
	}

	if ismqyaml {

		log.Printf("merge-mqsc-files: transforming mq yaml '%s' into mqsc '%s'\n",
			mqyaml, mqyamlout)

		err := mqsc.Outputmqsc(mqyaml, mqyamlout)
		if err != nil {
			return err
		}
	}

	// merge mqsc transform output with mqsc-ini file
	mqscini := _mqscini
	mqscic := _mqscic

	// check if mqsc-ini file exists
	ismqscini := true

	_, err = os.Stat(mqscini)
	if err != nil && os.IsNotExist(err) {
		log.Printf("merge-mqsc-files: mqsc ini file '%s' not found\n", mqscini)
		ismqscini = false

	} else if err != nil {
		return err
	}

	log.Printf("merge-mqsc-files: generating mqsc file '%s'\n", mqscic)

	err = MergeFiles(mqyamlout, ismqyaml, mqscini, ismqscini, mqscic)
	if err != nil {
		return err
	}

	return nil
}

func MergeFiles(mqyamlout string, ismqyaml bool, mqscini string, ismqscini bool, mqscic string ) error {

	out := "*\n"

	if ismqyaml {
		data, err := ioutil.ReadFile(mqyamlout)
		if err != nil {
			return err
		}
		out += string(data)
	}

	out += "*\n"

	if ismqscini {
		data, err := ioutil.ReadFile(mqscini)
		if err != nil {
			return err
		}

		out += string(data)
	}

	err := ioutil.WriteFile(mqscic, []byte(out), 0777)
	if err != nil {
		return err
	}

	return nil
}
