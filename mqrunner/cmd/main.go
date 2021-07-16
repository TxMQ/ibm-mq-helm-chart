package main

import (
	"fmt"
	"log"
	"os"
	"szesto.com/mqrunner/mqsc"
	"szesto.com/mqrunner/util"
	"szesto.com/mqrunner/webmq"
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

	// start runner
	log.Printf("mq runner %s starting...\n", qmgr)
	ctl := util.StartRunner()
	<-ctl

	//defer util.StopRunner()

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

	// todo: stop/start: make sure qmgr completely stopped
	// this is an edge use case

	//// check if qmgr is running
	//running, err := util.IsQmgrRunning(qmgr)
	//if err != nil {
	//	// log and exit
	//	log.Fatalf("is-qmgr-running: %v\n", err)
	//}
	//
	//// if running, stop qmgr
	//if running {
	//	err = util.StopQmgr(qmgr)
	//	if err != nil {
	//		// log and exit
	//		log.Fatalf("stop-qmgr: %v\n", err)
	//	}
	//
	//	log.Printf("qmgr %s stopped", qmgr)
	//}

	// start qeueue manager
	err = util.StartQmgr(qmgr)
	if err != nil {
		// log and exit
		log.Fatalf("start-qmgr: %v\n", err)

		// queue manager stop may take time
		// 2021/06/25 20:45:57 start-qmgr: IBM MQ queue manager 'qm10' ending.
	}

	log.Printf("qmgr %s started", qmgr)

	// set qmgr tls key repository
	if enabletls == "true" || enabletls == "1" {
		err = util.SetQmgrKeyRepoLocation(qmgr)
		if err != nil {
			// log and exit
			log.Fatalf("set-qmgr-key-repo-location: %v\n", err)
		}
	}

	// transform mq config yaml into mqsc commands
	mqconfigyaml := "/etc/mqm/mqsc/mqsc.yaml"
	startupmqsc := "/etc/mqm/startup.mqsc"

	log.Printf("transform mq config yaml '%s' into startup mqsc script '%s'\n",
		mqconfigyaml, startupmqsc)

	err = mqsc.Outputmqsc(mqconfigyaml, startupmqsc)
	if err != nil && os.IsNotExist(err) {
		// no mqconfig yaml file
		log.Printf("mq conifg yaml file '%s' not found\n", mqconfigyaml)

	} else if err != nil {
		// log and exit
		log.Fatalf("configure-webconsole: %v\n", err)

	} else {
		// apply mqsc commands
		log.Printf("applying '%s' mqsc file", startupmqsc)

		cout, err := util.RunmqscFromFile(qmgr, startupmqsc)
		if err != nil {
			cerr := string(cout)
			if len(cerr) > 0 {
				log.Printf("run-mqsc-from-file: %s\n", cerr)
			} else {
				log.Printf("run-mqsc-from-file: %v\n", err)
			}
		}
	}

	// apply mqsc ini commands
	mqscini := "/etc/mqm/mqsc/mqscini.mqsc"

	log.Printf("applying '%s' mqsc file", mqscini)

	cout, err := util.RunmqscFromFile(qmgr, mqscini)
	if err != nil {
		cerr := string(cout)
		if len(cerr) > 0 {
			log.Printf("run-mqsc-from-file: %s\n", cerr)
		} else {
			log.Printf("run-mqsc-from-file: %v\n", err)
		}
	}

	// configure webconsole
	log.Printf("%s\n", "configuring webconsole")

	err = webmq.ConfigureWebconsole()
	if err != nil {
		// log and exit
		log.Fatalf("configure-webconsole: %v\n", err)
	}

	// start webconsole
	log.Printf("%s\n", "starting mq web console")

	err = util.StartMqweb()
	if err != nil {
		// log and exit
		log.Fatalf("start-mq-web: %v\n", err)
	}

	log.Printf("mq runner %s running...\n", qmgr)

	<-ctl
	log.Printf("mq runner %s exiting...\n", qmgr)
}

func main() {
	Runmain()
}
