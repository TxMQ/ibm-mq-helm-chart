package main

import (
	"fmt"
	"os"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/mqrunner"
	"szesto.com/mqrunner/mqwebc"
	"szesto.com/mqrunner/perfmon"
	"szesto.com/mqrunner/probe"
	"szesto.com/mqrunner/qmgr"
	"szesto.com/mqrunner/util"
	"time"
)

func main() {
	startup := time.Now()

	qmname := os.Getenv("MQ_QMGR_NAME")
	logger.Logmsg(fmt.Sprintf("queue manager '%s' starting in startup role '%s'", qmname, qmgr.StartupRole()))

	//create /var/mqm directories
	if err := qmgr.CreateDirectories(); err != nil {
		logger.Logmsg(err)
	}

	logger.Runlogger()
	probe.StartProbes(qmname)
	mqrunner.StartMqrunner(qmname)
	mqrunner.WaitForRunnerReady()

	// config files are merged into local /etc/mqm directory
	if err := util.MergeMqscFiles(); err != nil {
		logger.Logmsg(err)
	}

	if err := util.MergeGitConfigFiles(); err != nil {
		logger.Logmsg(err)
	}

	maxtries := 5
	created := false
	for retryCount := 0; retryCount < maxtries && !created; retryCount++ {

		if err := qmgr.CreateQmgr(qmname, retryCount); err != nil {
			logger.Logmsg(err)

			if util.IsQmgrIniMissing(qmname, err) {
				// multi-instance leader did not create qmgr yet
				if err := qmgr.WaitForQmgrCreate(qmname); err != nil {
					logger.Logmsg(err)
				}

			} else {
				// create error, do not retry
				break
			}

		} else {
			created = true
		}
	}

	// start qmgr
	if err := qmgr.StartQmgr(qmname); err != nil {
		logger.Logmsg(err)
	}

	// tail logs
	qmgr.TailLogs(qmname)

	// start web console, /var/mqm/web directory
	mqwebc.StartWebconsole()

	// start state monitor
	qmgr.StartMonitor(qmname)

	// start perf-monitor
	perfmon.StartPerfMonitor()

	// running role (active, standby)
	if status, err := util.QmgrStatus(qmname, false); err != nil {
		logger.Logmsg(err)

	} else {
		switch status {
		case util.QmgrStatusEnumRunning():
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'active'", qmname))

			//import qmgr keystore
			t := time.Now()
			if _, _, err := qmgr.ImportQmgrKeystore(qmname); err != nil {
				logger.Logmsg(err)
			} else {
				_ = util.RefreshSsl(qmname)
			}
			logger.Logmsg(fmt.Sprintf("time to import qmgr keystore: %v", time.Since(t)))

			// cat autoconfig file

		case util.QmgrStatusEnumStandby():
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'standby'", qmname))

		case util.QmgrStatusEnumElsewhere():
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'running elsewhere', exiting", qmname))
			os.Exit(1)

		default:
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role (active|standby) is 'unknown'", qmname))
		}
	}

	logger.Logmsg(fmt.Sprintf("startup time: %v", time.Since(startup)))
	mqrunner.WaitForRunnerStop()

	logger.Logmsg("exiting")
}
