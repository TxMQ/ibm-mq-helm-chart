package main

import (
	"fmt"
	"log"
	"os"
	"szesto.com/mqrunner/util"
	"szesto.com/mqrunner/webmq"
	"time"
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

	if debug {
		_ = util.ListDir("/var/mqm")
		time.Sleep(1 * time.Second)
	}

	// create runtime directories
	err := util.CreateDirectories()
	if err != nil {
		// ignore for now
		fmt.Printf("create-directories: %v\n", err)
	}

	log.Printf("%s\n", "mq directories created")

	if debug {
		_ = util.ListDir("/var/mqm")
	}

	// get queue manager name
	qmgr := os.Getenv("MQ_QMGR_NAME")

	if debug {
		log.Printf("run-main: qmgr name: '%s'\n", qmgr)
	}

	// get qmgr log format basic|json

	// merge mqsc startup files
	err = util.MergeMqscFiles()
	if err != nil {
		log.Fatalf("merge-mqsc-files: %v\n", err)
	}

	// fetch and merge startup config files
	giturl := os.Getenv("GIT_CONFIG_URL")
	gitref := os.Getenv("GIT_CONFIG_REF")
	gitdir := os.Getenv("GIT_CONFIG_DIR")

	if debug {
		log.Printf("run-main: giturl '%s', gitref '%s', gitdir '%s'\n", giturl, gitref, gitdir)
	}

	if len(giturl) > 0 {

		fetchconf := util.FetchConfig{
			Url:           giturl,
			ReferenceName: gitref,
			Dir:           gitdir,
		}

		err = util.MergeGitConfigFiles(fetchconf)
		if err != nil {
			log.Fatalf("fetch-merge-config-files: %v\n", err)
		}
	}

	// start runner
	log.Printf("mq runner '%s' starting...\n", qmgr)
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
		if debug {
			log.Printf("run-main: qmgr %s does not exist, will be created\n", qmgr)
		}

		// create queue manager
		err = util.CreateQmgr(qmgr)
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

	// check mqsc syntax errors
	ok, err := util.CheckMqscSyntax(qmgr)
	if err != nil {
		log.Printf("%v\n", err)
	}

	if ok {
		if debug {
			log.Printf("run-main: %s\n", "startup mqsc syntax check passed")
		}
	} else {
		log.Printf("run-main: %s\n","startup mqsc commands contain syntax errors")
	}

	// tail system log: /var/mqm/errors/AMQERR01.LOG
	util.TailMqLog()

	// tail qmgr log: /var/mqm/qmgrs/{qmgr}/errors/AMQERR01.LOG
	util.TailQmgrLog(qmgr)

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

	// using default key repository
	// set qmgr tls key repository
	//if util.IsEnableTls() {
	//	err = util.SetQmgrKeyRepoLocation(qmgr)
	//	if err != nil {
	//		// log and exit
	//		log.Fatalf("set-qmgr-key-repo-location: %v\n", err)
	//	}
	//}

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
				// log and exit
				log.Fatalf("start-mq-web: %v\n", err)
			}
		}

	} else {
		log.Printf("%s\n", "webconsole is off")
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
