package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"szesto.com/mqrunner/logger"
)

func GetDebugFlag() bool {
	if debugenv := os.Getenv("MQRUNNER_DEBUG"); len(debugenv) > 0 {
		if debugenv == "1" || debugenv == "true" {
			return true
		}
	}
	return false
}

func ShowMounts() error {

	if GetDebugFlag() {
		log.Printf("show-mounts: cat /proc/mounts")
	}

	out, err := runcmd("cat", "/proc/mounts")
	if err != nil {
		fmt.Printf("show-mounts: %v\n", err)
		return err
	}

	fmt.Printf("show-mounts: %s\n", out)
	return nil
}

func ListDir(dir string) error {

	if GetDebugFlag() {
		log.Printf("list-dir: listing directory %s\n", dir)
	}

	out, err := runcmd("ls", "-l", dir)
	if err != nil {
		fmt.Printf("list-dir: out: %s, err: %v\n", out, err)
		return err
	}

	fmt.Printf("list-dir: %s\n", out)
	return nil
}

func runcmd(cmd string, args ...string) (string, error) {
	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	}

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		if len(string(out)) > 0 {
			cerr := string(out)
			return "", fmt.Errorf("%s%v", cerr, err)
		} else {
			return "", err
		}
	}

	cout := string(out)
	return cout, nil
}