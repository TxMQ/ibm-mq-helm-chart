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

	err = os.Setenv("MQ_ENABLE_TLS_NO_VAULT", "1")
	if err != nil {
		return
	}

	Runmain()
}
