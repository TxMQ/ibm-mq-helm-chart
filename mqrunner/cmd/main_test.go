package main

import (
	"os"
	"testing"
)

func TestRunMain(t *testing.T) {
	err := os.Setenv("MQ_QMGR_NAME", "qm10")
	if err != nil {
		return
	}

	Runmain()
}
