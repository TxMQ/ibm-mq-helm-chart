package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/mqmodel"
)

const _qmgrstarting = "starting"
const _qmgrrunning = "running"
const _qmrunningstandby = "runningstandby"
const _qmrunningelsewhere = "runningelsewhere"
const _qmgrnotrunning = "notrunning"
const _qmgrnotknown = "notknown"
const _qmgrstatusnotavailable = "statusnotavailable"

func QmgrStatusEnumRunning() string {
	return _qmgrrunning
}

func QmgrStatusEnumStandby() string {
	return _qmrunningstandby
}

func QmgrStatusEnumElsewhere() string {
	return _qmrunningelsewhere
}

func isEnvTrueValue(envvar string) bool {
	if value := os.Getenv(envvar); len(value) > 0 && (strings.ToLower(value) == "true" || value == "1") {
		return true
	}
	return false
}

func getEnv(envvar string) (bool, string) {
	if value := os.Getenv(envvar); len(value) > 0 {
		return true, value
	}
	return false, ""
}

func IsMultiInstance1() bool {
	// kubernetes
	if isEnvTrueValue("MULTI_INSTANCE_QMGR") {
		if isset, hostname := getEnv("HOSTNAME"); isset {
			if strings.HasSuffix(hostname, "-0") {
				return true
			} else if strings.HasSuffix(hostname, "-1") {
				return false
			}
		}
	}

	// plain docker
	if isEnvTrueValue("MULTI_INSTANCE_QMGR_1") {
		return true
	}

	return false
}

func IsMultiInstance2() bool {
	// kubernetes
	if isEnvTrueValue("MULTI_INSTANCE_QMGR") {
		if isset, hostname := getEnv("HOSTNAME"); isset {
			if strings.HasSuffix(hostname, "-0") {
				return false
			} else if strings.HasSuffix(hostname, "-1") {
				return true
			}
		}
	}

	// plain docker
	if isEnvTrueValue("MULTI_INSTANCE_QMGR_2") {
		return true
	}

	return false
}

func ApplyStartupConfig(qmgr string) error {

	cmdfile := GetMqscic()

	_, err := os.Stat(cmdfile)
	if err != nil && os.IsNotExist(err) {
		logger.Logmsg(fmt.Sprintf("command file '%s' does not exist", cmdfile))
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
		logger.Logmsg(fmt.Sprintf("run mqsc commands from '%s'\n", cmdfile))
	}

	out, err := RunmqscFromFile(qmgr, cmdfile)
	if err != nil {
		return err
	}

	if len(out) > 0 {
		logger.Logmsg(out)
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
	if IsMultiInstance1() {
		return createQmgrCmd(qmgr, icignore)

	} else if IsMultiInstance2() {
		return addMqinfCmd(qmgr)

	} else {
		return createQmgrCmd(qmgr, icignore)
	}
}

func addMqinfCmd(qmgr string) error {
	// dspmqinf -o command qm1
	// addmqinf -s QueueManager -v Name=qm1 -v Directory=qm1 -v Prefix=/var/mqm -v DataPath=/var/md/qm1

	args := []string {"-s", "QueueManager"}
	args = append(args, "-v", fmt.Sprintf("Name=%s", qmgr))
	args = append(args, "-v", fmt.Sprintf("Directory=%s", qmgr))
	args = append(args, "-v", fmt.Sprintf("Prefix=%s", "/var/mqm"))
	args = append(args, "-v", fmt.Sprintf("DataPath=/var/md/%s", qmgr))

	addmqinf := "/opt/mqm/bin/addmqinf"

	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("running: %s %s", addmqinf, strings.Join(args, " ")))
	}

	out, err := exec.Command(addmqinf, args...).CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return fmt.Errorf("%s%v", out, err)
		} else {
			return err
		}

	} else if len(out) > 0 {
		cout := string(out)
		logger.Logmsg(cout)

		if isQmgrIniMissing(qmgr, cout) {
			return fmt.Errorf("%s", cout)
		}
	}

	return nil
}

func IsQmgrIniMissing(qmgr string, err error) bool {
	return isQmgrIniMissing(qmgr, fmt.Sprintf("%v", err))
}

func isQmgrIniMissing(qmgr, out string) bool {
	// replace possible new lines
	errs := strings.ReplaceAll(out, "\n", " ")

	// AMQ5208E: File '/var/md/qm1/qm.ini' missing.
	amq5208e := strings.HasPrefix(errs, "AMQ5208E")

	// an expected stanza in an ini file is missing or contains errors, exit status 71
	status71 := strings.HasSuffix(errs, "exit status 71")

	return amq5208e || status71
}

func createQmgrCmd(qmgr string, icignore bool) error {

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
	// -ld - log path /mnt/data/ld -> /var/ld
	// -md - qmgr data path, /var/mqm/qmgrs; /mnt/data/md -> /var/md
	// -oa group - (default) authorization mode

	args := []string {"-c", "queue manager"}

	args = append(args, "-lc")

	if ismqscic {
		args = append(args, "-ic", mqscic)
	}

	if isqmini {
		args = append(args, "-ii", qmini)
	}

	args = append(args, "-ld", "/var/ld")
	args = append(args, "-md", "/var/md")
	args = append(args, "-u", deadLetterQeueue)
	args = append(args, "-p", qmgrPort)
	args = append(args, "-q")

	args = append(args, qmgr)

	//out, err := exec.Command("/opt/mqm/bin/crtmqm", "-c", "queue manager", "-lc",
	//	"-ic", mqscic,
	//	"-u", deadLetterQeueue, "-p", qmgrPort, "-q", qmgr).CombinedOutput()

	if debug {
		logger.Logmsg(fmt.Sprintf("running command: /opt/mqm/bin/crtmqm %s", strings.Join(args, " ")))
	}

	out, err := exec.Command("/opt/mqm/bin/crtmqm", args...).CombinedOutput()

	if debug {
		if len(string(out)) > 0 {
			if err == nil {
				logger.Logmsg(fmt.Sprintf("%s", string(out)))
			} else {
				logger.Logmsg(fmt.Sprintf("%s%v", string(out), err))
			}

		} else {
			logger.Logmsg(err)
		}
	}

	if err != nil {
		if out != nil {
			cerr := string(out)
			return fmt.Errorf("%v", cerr)
		} else {
			return err
		}
	}

	return nil
}

func StopMqweb() error {
	out, err := runcmd("/opt/mqm/bin/endmqweb")
	if err != nil {
		return err

	} else if len(out) > 0 {
		logger.Logmsg(out)
	}
	return nil
}

func StartMqweb() error {
	out, err := runcmd("/opt/mqm/bin/strmqweb")
	if err != nil {
		return err

	} else if len(out) > 0 {
		logger.Logmsg(out)
	}
	return nil
}

func StartQmgr(qmgr string) error {

	var out []byte
	var err error

	if IsMultiInstance1() || IsMultiInstance2() {
		// add -x argument for the multi-instance start
		out, err = exec.Command("/opt/mqm/bin/strmqm", "-x", qmgr).CombinedOutput()
	} else {
		// start queue manager
		out, err = exec.Command("/opt/mqm/bin/strmqm", qmgr).CombinedOutput()
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

func StopQmgr(qmgr string) error {

	var out []byte
	var err error

	// stop queue manager
	if IsMultiInstance1() || IsMultiInstance2() {
		out, err = exec.Command("/opt/mqm/bin/endmqm", "-x", qmgr).CombinedOutput()
	} else {
		out, err = exec.Command("/opt/mqm/bin/endmqm", qmgr).CombinedOutput()
	}

	if err != nil && len(out) > 0 {
		return fmt.Errorf("%s%v\n", string(out), err)
	} else if err != nil {
		return err
	}

	return nil
}

func IsQmgrRunning(qmgr string, silent bool) (bool, error) {

	st, err := QmgrStatus(qmgr, silent)
	if err != nil {
		return false, err
	}

	if st == _qmgrrunning {
		return true, nil
	}

	return false,nil
}

func IsQmgrRunningStandby(qmgr string, silent bool) (bool, error) {
	st, err := QmgrStatus(qmgr, silent)
	if err != nil {
		return false, err
	}

	if st == _qmrunningstandby {
		return true, nil
	}

	return false,nil
}

func QmgrConf(qmgr string) (bool, string, error) {

	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("checking if qmgr '%s' already configured", qmgr))
	}

	cout, err := runcmd("/opt/mqm/bin/dspmqinf", "-s", "QueueManager", qmgr)
	if err != nil {
		return false, "", err
	}

	return true, cout, nil
}

func QmgrExists(qmgr string) (bool, error) {

	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("check if queue manager %s exists", qmgr))
	}

	st, err := QmgrStatus(qmgr, false)
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

func parseQmgrStatusValue(cout string) string {
	// QMNAME(qm) STATUS(Running)
	// QMNAME(qm) STATUS(Running as standby)
	// QMNAME(qm) STATUS(Running elsewhere)
	// QMNAME(qm) STATUS(Starting)
	// QMNAME(qm) STATUS(Ended normally|immediately|unexpectedly)
	// QMNAME(qm) STATUS(Status not available)

	if ok, status := parseParenValue(cout, "STATUS"); ok {
		switch strings.ToLower(status) {
		case "running": return _qmgrrunning
		case "running as standby": return _qmrunningstandby
		case "running elsewhere": return _qmrunningelsewhere
		case "starting": return _qmgrstarting
		case "status not available": return _qmgrstatusnotavailable
		case "ended normally": return _qmgrnotrunning
		case "ended immediately": return _qmgrnotrunning
		case "ended unexpectedly": return _qmgrnotrunning
		default:
			logger.Logmsg(cout)
			return _qmgrstatusnotavailable
		}

	} else {
		return _qmgrstatusnotavailable
	}
}

func QmgrStatus(qmgr string, silent bool) (string, error) {
	debug := GetDebugFlag()

	args := []string{"-m", qmgr}

	if IsMultiInstance1() || IsMultiInstance2() {
		args = append(args, "-x")
	}

	dspmq := "/opt/mqm/bin/dspmq"

	if debug && !silent {
		logger.Logmsg(fmt.Sprintf("running: %s %s", dspmq, strings.Join(args, " ")))
	}

	out, err := exec.Command(dspmq, args...).CombinedOutput()

	if err != nil && len(out) > 0 {
		cerr := strings.TrimSpace(fmt.Sprintf("%v", err))
		cout := strings.TrimSpace(string(out))

		// AMQ7048E: The queue manager name is either not valid or not known.
		if strings.HasPrefix(cout, "AMQ7048E") {
			return _qmgrnotknown, nil

		} else if strings.HasPrefix(cerr,"wait" ) {
			// sometimes valid status is reported with an error: waitid|wait: no child processes
			if strings.HasPrefix(cout, "QMNAME") {
				return parseQmgrStatusValue(cout), nil
			}

		} else {
			logger.Logmsg(fmt.Sprintf("%s%v", out, err))
			return _qmgrstatusnotavailable, err
		}

	} else if err != nil {
		logger.Logmsg(err)
		return _qmgrstatusnotavailable, err
	}

	if debug && !silent {
		logger.Logmsg(string(out))
	}

	qmstatus := _qmgrstatusnotavailable
	cout := strings.TrimSpace(string(out))

	if strings.HasPrefix(cout, "QMNAME") {
		qmstatus = parseQmgrStatusValue(cout)
	}

	return qmstatus, nil
}

func Runmqsc(qmgr, command string) (string, error) {
	logger.Logmsg(fmt.Sprintf("running mqsc command: %s", command))

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

func RefreshSsl(qmgr string) error {
	// refresh security type(ssl)
	out, err := Runmqsc(qmgr, "refresh security type(ssl)")
	if err != nil {
		return err
	} else {
		logger.Logmsg(out)
	}
	return nil
}

func SetSslKeyRepo(qmgr, sslkeyr string) error {
	return SetQmgrParam(qmgr,"SSLKEYR", fmt.Sprintf("'%s'", sslkeyr))
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
	_ = mqmodel.ClearLdapBindPasswordEnv()
	return true
}