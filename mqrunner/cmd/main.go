package main

import (
	"fmt"
	"os"
	"szesto.com/mqrunner/util"
)

func main() {

	// env variables set in the pod template
	// MQ_QMGR_NAME - queue manager name

	// create runtime directories
	err := util.CreateDirectories()
	if err != nil {
		// decide what to do
	}

	// import certs into the keystore
	qmgr := os.Getenv("MQ_MGR_NAME")

	err = util.ImportCertificates(qmgr)
	if err != nil {
		// decide to continue or not
	}

	// check if qmgr exists
	exists, err := util.QmgrExists(qmgr)
	if err != nil {
		// decide to continue or not
	}

	if exists == false {
		// create queue manager
		err = util.CreateQmgr(qmgr)
		if err != nil {
			// decide what to do
		}
	}

	// check if qmgr is running
	running, err := util.IsQmgrRunning(qmgr)
	if err != nil {
		// decide to continue or not
	}

	// if running, stop qmgr
	if running {
		err = util.StopQmgr(qmgr)
	}

	// start qeueue manager
	err = util.StartQmgr(qmgr)
	if err != nil {
		// decide what to do
	}

	// set qmgr tls key repository and label

	// start webconsole

	// start runner
	<- util.StartRunner()

	fmt.Println("mqrunner exiting...")
}
