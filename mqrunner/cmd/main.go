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

func prepareQueueManager(qmgr string) error {
	if util.IsMultiInstance2() {
		log.Printf("prepare-qmgr: qmgr '%s' multi-instance standby, skip prepare\n", qmgr)
		return nil
	} else if util.IsMultiInstance1() {
		return prepareQueueManagerActive(qmgr)
	} else {
		return prepareQueueManagerActive(qmgr)
	}
}

func prepareQueueManagerActive(qmgr string) error {

	debug := util.GetDebugFlag()

	if debug {
		log.Printf("prepare-qmgr: qmgr name: '%s'\n", qmgr)
	}

	if debug {
		log.Printf("prepare-qmgr: debug flag set")
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
	if err := util.CreateDirectories(); err != nil {
		fmt.Printf("create-directories: %v\n", err)

		return err
	}

	log.Printf("%s\n", "mq directories created")

	//if debug {
	//	_ = util.ListDir("/var/mqm")
	//}

	// todo: qmgr log format basic|json

	// merge mqsc startup files
	if err := util.MergeMqscFiles(); err != nil {
		log.Printf("merge-mqsc-files: %v\n", err)

		// ignore merge error, continue
	}

	// fetch and merge config files
	if err := util.MergeGitConfigFiles(); err != nil {
		log.Printf("fetch-merge-config-files: %v\n", err)

		// ignore merge error, continue
	}

	return nil
}

func createQueueManager(qmgr string) error {

	if util.IsMultiInstance2() {
		log.Printf("create-qmgr: %s multi-instance-2, skip create...\n", qmgr)
		return nil
	}

	debug := util.GetDebugFlag()

	// check if qmgr already configured
	qmconf, msg, err := util.QmgrConf(qmgr)
	if err != nil {
		// log error message, continue as if qmgr does not exist
		log.Printf("run-main: %v\n", err)

	} else {
		log.Printf("%s\n", msg)
	}

	if qmconf == false {
		if debug {
			log.Printf("run-main: queue manager %s does not exist, will be created\n", qmgr)
		}

		// create queue manager, ignore ic file
		if err = util.CreateQmgr(qmgr, true); err != nil {
			// log error
			log.Printf("create-qmgr: %v\n", err)

			log.Printf("queue manager %s create failed", qmgr)

			return err

		} else {
			log.Printf("queue manager %s created", qmgr)
		}

	} else {
		if debug {
			log.Printf("run-main: queue manager '%s' exists", qmgr)
		}
	}

	return nil
}

func postCreateQueueManager(qmgr string) error {

	debug := util.GetDebugFlag()

	// tail system log: /var/mqm/errors/AMQERR01.LOG
	util.TailMqLog()

	// tail qmgr log: /var/mqm/qmgrs/{qmgr}/errors/AMQERR01.LOG
	util.TailQmgrLog(qmgr)

	if util.IsEnableTls() {
		if util.IsMultiInstance2() {
			// skip cert import
			log.Printf("post-create-qmgr: '%s' multi-instance-2 skip tls cert import\n", qmgr)
		} else {
			// import certs into the keystore
			if err := util.ImportCertificates(qmgr); err != nil {
				// log error
				log.Printf("import-certificates: %v\n", err)

				// return error?
			}
		}

	} else {
		if debug {
			log.Printf("run-main: TLS is not enabled")
		}
	}

	// start qeueue manager
	if err := util.StartQmgr(qmgr); err != nil {
		// log error
		log.Printf("start-qmgr: %v\n", err)

		log.Printf("Queue manager %s did not start\n", qmgr)

		return err
	}

	log.Printf("Queue manager '%s' started", qmgr)

	// configure webconsole
	if isStartMqweb() || isConfigureMqweb() {
		log.Printf("%s\n", "configuring webconsole")

		err := webmq.ConfigureWebconsole()
		if err != nil {
			// log error
			log.Printf("configure-webconsole: %v\n", err)

			if isStartMqweb() {
				log.Printf("web console configuration failed, web console will not be started")
			}

		} else {

			if isStartMqweb() {
				log.Printf("%s\n", "starting mq web console")
				log.Printf("%s\n", "web console will connect to ldap, if taking too long check ldap server")

				err = util.StartMqweb()
				if err != nil {
					// log error
					log.Printf("start-mq-web: %v\n", err)
				}
			}
		}

	} else {
		log.Printf("%s\n", "webconsole is off")
	}

	// apply startup configuration
	if err := util.ApplyStartupConfig(qmgr); err != nil {
		// log error message
		log.Printf("run-main: %v\n", err)
	}

	return nil
}

func Runmain() {

	// get queue manager name
	qmgr := os.Getenv("MQ_QMGR_NAME")
	log.Printf("run-main: qmgr name: '%s'\n", qmgr)

	// prepare queue manager
	err := prepareQueueManager(qmgr)
	if err != nil {
		// log error
		log.Printf("run-main: %v\n", err)
	}

	// start mq runner
	log.Printf("mq runner '%s' starting...\n", qmgr)
	ctl := util.StartRunner(qmgr)
	<-ctl

	// create queue manager
	err = createQueueManager(qmgr)
	if err != nil {
		fmt.Printf("run-main: %v\n", err)

		fmt.Printf("run-main: creating queue manager failed, queue manager %s is not fully functional\n", qmgr)
	}

	// post-create queue manager
	err = postCreateQueueManager(qmgr)
	if err != nil {
		// log error
		fmt.Printf("run-main: %v\n", err)

		fmt.Printf("run-main: some post-create steps failed, queue manager %s is not fully functional\n", qmgr)
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
