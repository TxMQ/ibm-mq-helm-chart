package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const _qmgrrunning = "running"
const _qmgrnotrunning = "notrunning"
const _qmgrnotknown = "notknown"

func CreateDirectories() error {

	//
	// create mq directories
	// this command requires su on rpm install
	//
	out, err := exec.Command("/opt/mqm/bin/crtmqdir", "-f", "-a").CombinedOutput()
	if err != nil {
		if out != nil {
			cerr := string(out)
			return fmt.Errorf("%v\n", cerr)
		} else {
			return err
		}
	}

	return nil
}

func CreateQmgr(qmgr string) error {

	// create queue manager
	//
	// crtmqm -c "queue manager" -ic mqsi-file-path -ii ini-file-path -lc -p 1414 -q -u SYSTEM.DEAD.LETTER.QUEUE
	// -lc - circular logging
	// -ii argument is not passed yet
	// -ic argument is not passed yet
	// -md - qmgr data path, /var/mqm/qmgrs
	// -oa group - (default) authorization mode
	out, err := exec.Command("/opt/mqm/bin/crtmqm", "-c", "queue manager", "-lc",
		"-u", "SYSTEM.DEAD.LETTER.QUEUE", "-p", "1414", "-q", qmgr).CombinedOutput()

	if err != nil {
		if out != nil {
			cerr := string(out)
			return fmt.Errorf("%v\n", cerr)
		} else {
			return err
		}
	}
	return nil
}

func StartQmgr(qmgr string) error {

	// start queue manager
	out, err := exec.Command("/opt/mqm/bin/strmqm", qmgr).CombinedOutput()

	if err != nil {
		if out != nil {
			cerr := string(out)
			return fmt.Errorf("%v\n", cerr)
		} else {
			return err
		}
	}

	return nil
}

func StopQmgr(qmgr string) error {
	// stop queue manager
	out, err := exec.Command("/opt/mqm/bin/endmqm", qmgr).CombinedOutput()

	if err != nil {
		if out != nil {
			cerr := string(out)
			return fmt.Errorf("%v\n", cerr)
		} else {
			return err
		}
	}

	return nil
}

func IsQmgrRunning(qmgr string) (bool, error) {

	st, err := QmgrStatus(qmgr)
	if err != nil {
		return false, err
	}

	if st == _qmgrrunning {
		return true, nil
	}

	return false,nil
}

func QmgrExists(qmgr string) (bool, error) {

	st, err := QmgrStatus(qmgr)
	if err != nil {
		return false, err
	}

	if st == _qmgrnotknown {
		return false, nil
	}

	return true,nil
}

func parseParenValue(input, keyword string) (bool, string) {
	n1 := strings.Index(input, keyword)
	if n1 < 0 { return false, "" }

	n11 := strings.Index(input[n1:], "(")
	if n11 < 0 { return false, "" }

	n12 := strings.Index(input[n1:], ")")
	if n12 < 0 { return false, "" }

	value := input[n1+n11+1:n1+n12]
	return true, value
}

func QmgrStatus(qmgr string) (string, error) {

	out, err := exec.Command("/opt/mqm/bin/dspmq", "-m", qmgr).CombinedOutput()

	if err != nil {
		cerr := strings.TrimSpace(string(out))

		// AMQ7048E: The queue manager name is either not valid or not known.
		if strings.HasPrefix(cerr, "AMQ7048E") {
			return _qmgrnotknown, nil
		}

		return "", err
	}

	cout := strings.TrimSpace(string(out))

	// QMNAME(qm) STATUS(Running)
	if strings.HasPrefix(cout, "QMNAME") {

		ok, status := parseParenValue(cout, "STATUS")

		if ok && strings.ToLower(status) == "running" {
			return _qmgrrunning, nil
		}
	}

	// QMNAME(qm)  STATUS(Ended normally|immediately)
	return _qmgrnotrunning, nil
}

func Runmqsc(qmgr, command string) (string, error) {
	var cmdfile = filepath.Join(os.TempDir(), "cmd.mqsc")

	err := ioutil.WriteFile(cmdfile, []byte(command), 0777)
	if err != nil {
		return "", err
	}

	return RunmqscFromFile(qmgr, cmdfile)
}

func RunmqscFromFile(qmgr, cmdfile string) (string, error) {

	//var cmdfile = "/tmp/cmd.mqsc"
	//var cmdfile = filepath.Join(os.TempDir(), "cmd.mqsc")
	//
	//err := ioutil.WriteFile(cmdfile, []byte(command), 0777)
	//if err != nil {
	//	return "", err
	//}

	// see if command file exists
	_, err := os.Stat(cmdfile)
	if err != nil {
		return "", err
	}

	cout, err := exec.Command("/opt/mqm/bin/runmqsc", "-e", "-f", cmdfile, qmgr).CombinedOutput()
	if err != nil {
		if cout != nil {
			// AMQ8118E: IBM MQ queue manager does not exist.
			cerr := string(cout)
			if idx := strings.Index(cerr, "AMQ8118E"); idx >= 0 {
				return "", fmt.Errorf("AMQ8118E: IBM MQ queue manager %s does not exist", qmgr)
			} else {
				return cerr, err
			}
		} else {
			return "", err
		}
	}

	return strings.TrimSpace(string(cout)), nil
}

func SetQmgrParam(qmgr, param, value string) error {

	paramuc := strings.TrimSpace(strings.ToUpper(param))
	cmd := fmt.Sprintf("alter qmgr %s(%s)", paramuc, strings.TrimSpace(value))

	_, err := Runmqsc(qmgr, cmd)
	if err != nil {
		return err
	}

	return nil
}

func GetQmgrParam(qmgr, param string) (string, error) {

	paramuc := strings.TrimSpace(strings.ToUpper(param))
	cmd := fmt.Sprintf("display qmgr %s", paramuc)

	out, err := Runmqsc(qmgr, cmd)
	if err != nil {
		return "", err
	}

	if ok, value := parseParenValue(out, paramuc); ok {
		return value, nil
	}

	return "", fmt.Errorf("qmgr parameter %s not found", param)
}

func RefreshSsl() error {
	// todo
	return nil
}

func SetSslKeyRepo(qmgr, sslkeyr string) error {
	return SetQmgrParam(qmgr,"SSLKEYR", sslkeyr)
}

func GetSslKeyRepo(qmgr string) (string, error) {
	return GetQmgrParam(qmgr, "SSLKEYR")
}

func SetCertLabel(qmgr, certlabel string) error {
	return SetQmgrParam(qmgr,"CERTLABL", certlabel)
}

func GetCertLabel(qmgr string) (string, error) {
	return GetQmgrParam(qmgr, "CERTLABL")
}
