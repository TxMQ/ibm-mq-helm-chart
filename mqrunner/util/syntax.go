package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CheckMqscSyntax(qmgr string) (bool, error) {

	cmdfile := GetMqscic()

	log.Printf("check-mqsc-syntax: checking file '%s'", cmdfile)

	_, err := os.Stat(cmdfile)
	if err != nil && os.IsNotExist(err) {
		if GetDebugFlag() {
			fmt.Printf("check-mqsc-syntax: file '%s' does not exist\n", cmdfile)
		}
		return true, nil

	} else if err != nil {
		return false, err
	}

	cout, err := exec.Command("/opt/mqm/bin/runmqsc", "-e", "-v", "-f", cmdfile, qmgr).CombinedOutput()
	if err != nil {
		if cout != nil {
			// AMQ8118E: IBM MQ queue manager does not exist.
			cerr := string(cout)
			if idx := strings.Index(cerr, "AMQ8118E"); idx >= 0 {
				return false, fmt.Errorf("AMQ8118E: IBM MQ queue manager %s does not exist", qmgr)
			} else {
				return false, fmt.Errorf("out: %s, err: %v\n", string(cout), err)
			}
		} else {
			return false, err
		}
	}

	syntaxmsg := strings.TrimSpace(string(cout))

	log.Printf("check-mqsc-syntax: %s\n", syntaxmsg)

	// look for: No commands have a syntax error
	const noerr = "No commands have a syntax error"

	n := strings.Index(strings.TrimSpace(string(cout)), noerr)
	if n >= 0 {
		log.Printf("check-mqsc-syntax: file %s: %s", cmdfile, noerr)
		return true, nil

	} else {
		log.Printf("check-mqsc-syntax: file %s contains syntax errors", cmdfile)

		// rename mqsc file
		badfile := fmt.Sprintf("%s.badsyntax", cmdfile)

		err = os.Rename(cmdfile, badfile)
		if err != nil {
			return false, err
		}

		log.Printf("check-mqsc-syntax: mqsc file %s renamed to %s\n", cmdfile, badfile)

		return false, nil
	}
}
