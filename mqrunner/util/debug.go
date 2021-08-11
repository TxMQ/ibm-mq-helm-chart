package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		fmt.Printf("list-dir: out: %s, err: %v\n", string(out), err)
		return err
	}

	fmt.Printf("list-dir: %s\n", out)
	return nil
}

func runcmd(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)

	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out

	if err := c.Run(); err != nil {
		return "", err
	}

	return c.String(), nil
}