package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"szesto.com/mqrunner/mqsc"
)

const _qmgrrunning = "running"
const _qmgrnotrunning = "notrunning"
const _qmgrnotknown = "notknown"

func ApplyStartupConfig(qmgr string) error {

	cmdfile := GetMqscic()

	_, err := os.Stat(cmdfile)
	if err != nil && os.IsNotExist(err) {
		log.Printf("apply-startup-config: command file '%s' does not exist\n", cmdfile)
		return nil

	} else if err != nil {
		return err
	}

	// check mqsc syntax errors
	//ok, err := CheckMqscSyntax(qmgr, cmdfile)
	//if err != nil {
	//	log.Printf("%v\n", err)
	//}
	//
	//if !ok {
	//	return fmt.Errorf("%s\n","startup mqsc commands contain syntax errors")
	//}
	//
	//if debug {
	//	log.Printf("apply-startup-config: %s\n", "startup mqsc syntax check passed")
	//}

	if GetDebugFlag() {
		log.Printf("apply-startup-config: run mqsc commands from '%s'\n", cmdfile)
	}

	out, err := RunmqscFromFile(qmgr, cmdfile)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		log.Printf("apply-startup-config: %s\n", out)
	}

	// parse out syntax message

	return nil
}

func CreateDirectories() error {

	//
	// create mq directories
	// this command requires su on rpm install
	//
	if GetDebugFlag() {
		log.Printf("create-directories: %s\n","/opt/mqm/bin/crtmqdir_setuid -f -a")
	}

	out, err := exec.Command("/opt/mqm/bin/crtmqdir_setuid", "-f", "-a").CombinedOutput()
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

func DeleteQmgr(qmgr string) error {

	out, err := runcmd("/opt/mqm/bin/dltmqm", qmgr)
	if err != nil {
		return err
	}

	log.Printf("delete-qmgr: %s\n", out)
	return nil
}

func CreateQmgr(qmgr string, icignore bool) error {

	debug := GetDebugFlag()

	// qmgr parameters - may change
	qmgrPort := "1414"
	deadLetterQeueue := "SYSTEM.DEAD.LETTER.QUEUE" // todo: env variable
	mqscic := GetMqscic()
	qmini := GetQmini()

	ismqscic := true
	isqmini := true

	_, err := os.Stat(mqscic)
	if err != nil || icignore {
		ismqscic = false
	}

	_, err = os.Stat(qmini)
	if err != nil {
		isqmini = false
	}

	// create queue manager
	//
	// crtmqm -c "queue manager" -ic mqsc-file-path -ii ini-file-path -lc -p 1414 -q -u SYSTEM.DEAD.LETTER.QUEUE
	// -lc - circular logging
	// -ii qm.ini file
	// -ic mqsc file
	// -md - qmgr data path, /var/mqm/qmgrs
	// -oa group - (default) authorization mode

	args := []string {"-c", "queue manager"}

	args = append(args, "-lc")

	if ismqscic {
		args = append(args, "-ic", mqscic)
	}

	if isqmini {
		args = append(args, "-ii", qmini)
	}

	args = append(args, "-u", deadLetterQeueue)
	args = append(args, "-p", qmgrPort)
	args = append(args, "-q")

	args = append(args, qmgr)

	//out, err := exec.Command("/opt/mqm/bin/crtmqm", "-c", "queue manager", "-lc",
	//	"-ic", mqscic,
	//	"-u", deadLetterQeueue, "-p", qmgrPort, "-q", qmgr).CombinedOutput()

	if debug {
		log.Printf("create-qmgr: running command: /opt/mqm/bin/crtmqm %s\n", strings.Join(args, " "))
	}

	out, err := exec.Command("/opt/mqm/bin/crtmqm", args...).CombinedOutput()

	if debug {
		if len(string(out)) > 0 {
			log.Printf("create-qmgr: out: %s, err: %v\n", string(out), err)
		} else {
			log.Printf("create-qmgr: err: %v\n", err)
		}
	}

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

func StartMqweb() error {

	// start mq web console
	out, err := exec.Command("/opt/mqm/bin/strmqweb").CombinedOutput()

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

func QmgrConf(qmgr string) (bool, string, error) {

	if GetDebugFlag() {
		log.Printf("qmgr-conf: check if qmgr '%s' already configured\n", qmgr)
	}

	out, err := runcmd("/opt/mqm/bin/dspmqinf", "-s", "QueueManager", qmgr)
	if err != nil {
		if len(out) > 0 {
			cerr := string(out)
			return false, "", fmt.Errorf("out: %s, err: %v\n", cerr, err)
		} else {
			return false, "", err
		}
	}

	cout := string(out)
	return true, cout, nil
}

func QmgrExists(qmgr string) (bool, error) {

	if GetDebugFlag() {
		log.Printf("qmgr-exists: check if queue manager %s exists\n", qmgr)
	}

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

	debug := GetDebugFlag()

	if debug {
		log.Printf("qmgr-status: running command: /opt/mqm/bin/dspmq -m %s", qmgr)
	}

	out, err := exec.Command("/opt/mqm/bin/dspmq", "-m", qmgr).CombinedOutput()

	if debug {
		if len(string(out)) > 0 {
			log.Printf("qmgr-status: out: %s, err: %v\n", string(out), err)
		} else {
			log.Printf("qmgr-status: err: %v\n", err)
		}
	}

	if err != nil {
		cerr := strings.TrimSpace(string(out))

		// AMQ7048E: The queue manager name is either not valid or not known.
		if strings.HasPrefix(cerr, "AMQ7048E") {

			if debug {
				log.Printf("qmgr-status: qmgr %s status is '%s'\n", qmgr, _qmgrnotknown)
			}

			return _qmgrnotknown, nil
		}

		return "", err
	}

	cout := strings.TrimSpace(string(out))

	// QMNAME(qm) STATUS(Running)
	if strings.HasPrefix(cout, "QMNAME") {

		ok, status := parseParenValue(cout, "STATUS")

		if ok && strings.ToLower(status) == "running" {

			if debug {
				log.Printf("qmgr-status: qmgr %s status is '%s'\n", qmgr, _qmgrrunning)
			}

			return _qmgrrunning, nil
		}
	}

	// QMNAME(qm)  STATUS(Ended normally|immediately)
	if debug {
		log.Printf("qmgr-status: qmgr %s status is '%s'\n", qmgr, _qmgrnotrunning)
	}

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
				return "", fmt.Errorf("run-mqsc-from-file: %s, %v\n", cerr, err)
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

func ClearEnvSecrets() bool {
	_ = mqsc.ClearLdapBindPasswordEnv()
	return true
}