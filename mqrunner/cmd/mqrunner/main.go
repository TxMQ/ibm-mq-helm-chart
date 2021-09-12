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

	if err := qmgr.CreateQmgr(qmname); err != nil {
		logger.Logmsg(err)

		if util.IsQmgrIniMissing(qmname, err) {
			// multi-instance leader did not create qmgr
			if err := qmgr.WaitForQmgrCreate(qmname); err != nil {
				logger.Logmsg(err)
			}
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
	tries := 0
	maxtries := 3
	isrunrole := false
	for tries < maxtries && !isrunrole {
		if qmgr.IsRunningRoleActive(qmname) {
			isrunrole = true
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'active'", qmname))

			//import qmgr keystore
			t := time.Now()
			if _, _, err := qmgr.ImportQmgrKeystore(qmname); err != nil {
				logger.Logmsg(err)
			}
			logger.Logmsg(fmt.Sprintf("time to import qmgr keystore: %v", time.Since(t)))

			// cat autoconfig file

		} else if qmgr.IsRunningRoleStandby(qmname) {
			isrunrole = true
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'standby'", qmname))

		} else {
			// uknown state, retry
			logger.Logmsg(fmt.Sprintf("qmgr '%s' running role (active|standby) is 'unknown'", qmname))
			tries++
			time.Sleep(5*time.Second)
		}
	}

	logger.Logmsg(fmt.Sprintf("startup time: %v", time.Since(startup)))
	mqrunner.WaitForRunnerStop()

	logger.Logmsg("exiting")
}
