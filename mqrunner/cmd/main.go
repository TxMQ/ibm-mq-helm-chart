package main

import (
	"fmt"
	"log"
	"os"
	"szesto.com/mqrunner/util"
)

func Runmain() {

	// env variables set in the pod template
	// MQ_QMGR_NAME - queue manager name

	// create runtime directories
	err := util.CreateDirectories()
	if err != nil {
		// ignore for now
		fmt.Printf("create-directories: %v\n", err)
	}

	log.Printf("%s\n", "mq directories created")

	// import certs into the keystore
	qmgr := os.Getenv("MQ_QMGR_NAME")

	enabletls := os.Getenv("MQ_ENABLE_TLS")
	if enabletls == "true" || enabletls == "1" {

		err = util.ImportCertificates(qmgr)
		if err != nil {
			// log and exit
			log.Fatalf("import-certificates: %v\n", err)
		}
	}

	// check if qmgr exists
	exists, err := util.QmgrExists(qmgr)
	if err != nil {
		// log and exit
		log.Fatalf("qmgr-exists: %v\n", err)
	}

	if exists == false {
		// create queue manager
		err = util.CreateQmgr(qmgr)
		if err != nil {
			// log and exit
			log.Fatalf("create-qmgr: %v\n", err)
		}

		log.Printf("qmgr %s created", qmgr)
	}

	// check if qmgr is running
	running, err := util.IsQmgrRunning(qmgr)
	if err != nil {
		// log and exit
		log.Fatalf("is-qmgr-running: %v\n", err)
	}

	// if running, stop qmgr
	if running {
		err = util.StopQmgr(qmgr)
		if err != nil {
			// log and exit
			log.Fatalf("stop-qmgr: %v\n", err)
		}

		// wait for the queue manager to stop

		log.Printf("qmgr %s stopped", qmgr)
	}

	// start runner
	log.Printf("mq runner %s starting...\n", qmgr)
	ctl := util.StartRunner()
	<-ctl

	// start qeueue manager
	err = util.StartQmgr(qmgr)
	if err != nil {
		// log and exit
		log.Fatalf("start-qmgr: %v\n", err)
	}

	log.Printf("qmgr %s started", qmgr)

	// set qmgr tls key repository and label

	// start webconsole

	fmt.Printf("mq runner %s running...\n", qmgr)

	<-ctl
	fmt.Printf("mq runner %s exiting...\n", qmgr)
}

func main() {
	Runmain()
}
