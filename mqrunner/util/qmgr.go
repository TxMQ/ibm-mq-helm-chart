package util

import (
	"fmt"
	"os/exec"
	"strings"
)

func CreateDirectories() error {

	// create mq directories
	cmd := exec.Command("/opt/mqm/bin/crtmqdir", "-f", "-a")
	err := cmd.Run()
	if err != nil {
		// complains about chmod 2775 on /mnt/mqm/data
		fmt.Printf("%v\n", err)
		return err
	}

	fmt.Println("directories created...")
	return nil
}

func CreateQmgr(qmgrname string) error {

	// todo
	// create queue manager
	cmd := exec.Command("/opt/mqm/bin/crtmqm", "-c", "qm", "-p", "1414", "-q", "qm")
	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		return nil
	}

	fmt.Println("queue manager created...")
	return nil
}

func StartQmgr(qmgrname string) error {

	// apply mq.ini overrides: -ii flag
	// apply mqsi-ini scripts - ic flag

	// todo
	// start queue manager
	cmd := exec.Command("/opt/mqm/bin/strmqm", "qm")
	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		return err
	}

	fmt.Println("queue manager started...")
	return nil
}

func StopQmgr(qmgrname string) error {
	// todo
	return nil
}

func IsQmgrRunning(qmgrname string) (bool, error) {

	st, err := QmgrStatus(qmgrname)
	if err != nil {
		return false, err
	}

	if st == "running" {
		return true, nil
	}

	return false,nil
}

func QmgrExists(qmgrname string) (bool, error) {

	st, err := QmgrStatus(qmgrname)
	if err != nil {
		return false, err
	}

	if st == "notknown" {
		return false, nil
	}

	return true,nil
}

func QmgrStatus(qmgrname string) (string, error) {

	out, err := exec.Command("/opt/mqm/bin/dspmq", "-m", qmgrname).CombinedOutput()
	if err != nil {
		return "", err
	}

	// AMQ7048E: The queue manager name is either not valid or not known.
	if strings.HasPrefix(strings.TrimSpace(string(out)), "AMQ7048E") {
		return "notknown", nil
	}

	// QMNAME(qm) STATUS(Running)
	if strings.HasPrefix(strings.TrimSpace(string(out)), "QMNAME") {

		var pname, pstatus string
		_, err := fmt.Scanf(strings.TrimSpace(string(out)), "QMNAME(%s) STATUS(%s)", &pname, &pstatus)
		if err != nil {
			return "", err
		}

		if strings.ToLower(pstatus) == "running" {
			return "running", nil
		}
	}

	// QMNAME(qm)  STATUS(Ended normally)
	return "notrunning", nil
}

func setQmgrParam(name, value string) error {

	cname := fmt.Sprintf("echo \"alter util %s(%s)\" | /opt/mqm/bin/runmqsc -e",
		strings.ToUpper(strings.TrimSpace(name)), strings.TrimSpace(value))

	err := exec.Command(cname).Run()
	return err
}

func getQmgrParam(param string) (string, error) {
	paramuc := strings.TrimSpace(strings.ToUpper(param))
	name := fmt.Sprintf("display util %s | /opt/mqm/bin/runmqsc -e | grep %s", paramuc, paramuc)

	cout, err := exec.Command(name).Output()
	if err != nil {
		return "", err
	}

	// qmname(qm) maxhands(256)
	// sslkeyr(/mnt/mqm/data/.../key)
	// qmname(qm)

	for _, p := range strings.Split(strings.TrimSpace(string(cout)), " ") {

		var paramName, paramValue string
		_, err = fmt.Scanf(p, "%s(%s)", &paramName, paramValue)
		if err != nil {
			return "", err
		}

		if paramName == paramuc {
			return paramValue, nil
		}
	}

	return "", fmt.Errorf("util parameter %s not found", param)
}

func refreshSsl() error {
	// todo
	return nil
}

func setSslKeyRepo(sslkeyr string) error {
	return setQmgrParam("SSLKEYR", sslkeyr)
}

func getSslKeyRepo() (string, error) {
	return getQmgrParam("SSLKEYR")
}

func setCertLabel(certlabel string) error {
	return setQmgrParam("CERTLABL", certlabel)
}

func getCertLabel() (string, error) {
	return getQmgrParam("CERTLABL")
}