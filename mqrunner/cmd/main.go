package main

import (
	"fmt"
	"log"
	"os"
	"szesto.com/mqrunner/util"
	"szesto.com/mqrunner/webmq"
)

func isStartMqweb() bool {
	return os.Getenv("MQ_START_MQWEB") == "1"
}

func isConfigureMqweb() bool {
	return os.Getenv("MQ_CONFIGURE_MQWEB") == "1"
}

func Runmain() {

	debug := util.GetDebugFlag()

	if debug {
		log.Printf("run-main: debug flag set")
	}

	//if debug {
	//	_ = util.ShowMounts()
	//}

	// env variables set in the pod template
	// MQ_QMGR_NAME - queue manager name

	//if debug {
	//	_ = util.ListDir("/var/mqm")
	//	time.Sleep(1 * time.Second)
	//}

	// create runtime directories
	err := util.CreateDirectories()
	if err != nil {
		// ignore for now
		fmt.Printf("create-directories: %v\n", err)
	}

	log.Printf("%s\n", "mq directories created")

	//if debug {
	//	_ = util.ListDir("/var/mqm")
	//}

	// get queue manager name
	qmgr := os.Getenv("MQ_QMGR_NAME")

	if debug {
		log.Printf("run-main: qmgr name: '%s'\n", qmgr)
	}

	// todo: qmgr log format basic|json

	// merge mqsc startup files
	err = util.MergeMqscFiles()
	if err != nil {
		log.Fatalf("merge-mqsc-files: %v\n", err)
	}

	// fetch and merge config files
	err = util.MergeGitConfigFiles()
	if err != nil {
		log.Printf("fetch-merge-config-files: %v\n", err)
	}

	// start runner
	log.Printf("mq runner '%s' starting...\n", qmgr)
	ctl := util.StartRunner()
	<-ctl

	//defer util.StopRunner()

	// check if qmgr already configured
	qmconf, msg, err := util.QmgrConf(qmgr)
	if err != nil {
		log.Printf("run-main: %v\n", err)
	} else {
		log.Printf("%s\n", msg)
	}

	if qmconf == false {
		if debug {
			log.Printf("run-main: qmgr %s does not exist, will be created\n", qmgr)
		}

		// create queue manager, ignore ic file
		err = util.CreateQmgr(qmgr, true)
		if err != nil {
			// log and exit
			log.Fatalf("create-qmgr: %v\n", err)
		}

		log.Printf("qmgr %s created", qmgr)

	} else {
		if debug {
			log.Printf("run-main: qmgr '%s' exists", qmgr)
		}
	}

	// tail system log: /var/mqm/errors/AMQERR01.LOG
	util.TailMqLog()

	// tail qmgr log: /var/mqm/qmgrs/{qmgr}/errors/AMQERR01.LOG
	util.TailQmgrLog(qmgr)

	if util.IsEnableTls() {
		// import certs into the keystore
		err = util.ImportCertificates(qmgr)
		if err != nil {
			// log and exit
			log.Fatalf("import-certificates: %v\n", err)
		}
	} else {
		if debug {
			log.Printf("run-main: TLS is not enabled")
		}
	}

	// start qeueue manager
	err = util.StartQmgr(qmgr)
	if err != nil {
		// log and exit
		log.Fatalf("start-qmgr: %v\n", err)

		// queue manager stop may take time
		// 2021/06/25 20:45:57 start-qmgr: IBM MQ queue manager 'qm10' ending.
	}

	log.Printf("qmgr '%s' started", qmgr)

	// configure webconsole
	if isStartMqweb() || isConfigureMqweb() {
		log.Printf("%s\n", "configuring webconsole")

		err = webmq.ConfigureWebconsole()
		if err != nil {
			// log and exit
			log.Fatalf("configure-webconsole: %v\n", err)
		}

		log.Printf("%s\n", "starting mq web console")

		if isStartMqweb() {
			err = util.StartMqweb()
			if err != nil {
				// log
				log.Printf("start-mq-web: %v\n", err)
			}
		}

	} else {
		log.Printf("%s\n", "webconsole is off")
	}

	// apply startup configuration
	if err = util.ApplyStartupConfig(qmgr); err != nil {
		log.Printf("run-main: %v\n", err)
	}

	// clear env var secrets
	util.ClearEnvSecrets()

	log.Printf("mq runner '%s' running...\n", qmgr)

	<-ctl
	log.Printf("mq runner '%s' exiting...\n", qmgr)
}

func main() {
	Runmain()
}
