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

	if qmgr.IsStartupLeader() {
		// create /var/mqm directories
		if err := qmgr.CreateDirectories(); err != nil {
			logger.Logmsg(err)
		}

		// create qmgr
		if err := qmgr.CreateQmgr(qmname); err != nil {
			logger.Logmsg(err)
		}

	} else {
		// wait for qmgr to be created by the leader
		_ = qmgr.WaitForQmgrCreate(qmname)
	}

	// tail logs
	qmgr.TailLogs(qmname)

	// start qmgr
	if err := qmgr.StartQmgr(qmname); err != nil {
		logger.Logmsg(err)
	}

	// start webc, web is starting async; web keystore is local file in /etc/mqm
	mqwebc.StartWebconsole()

	// let qmgr start...
	logger.Logmsg(fmt.Sprintf("pausing for %d seconds for qmgr '%s' to start", 5, qmname))
	time.Sleep(5 * time.Second)

	// get running role (active, standby)
	if qmgr.IsRunningRoleActive(qmname) {

		logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'active'", qmname))

		// todo: import qmgr keystore into local file system (not /var/mqm)
		if err := qmgr.ImportQmgrKeystore(qmname); err != nil {
			logger.Logmsg(err)
		}

		// cat autoconfig file

	} else if qmgr.IsRunningRoleStandby(qmname) {
		logger.Logmsg(fmt.Sprintf("qmgr '%s' running role is 'standby'", qmname))

	} else {
		// uknown state
		logger.Logmsg(fmt.Sprintf("qmgr '%s' running role (active|standby) is 'unknown'", qmname))
	}

	// start state monitor
	qmgr.StartMonitor(qmname)

	// start perf-monitor
	perfmon.StartPerfMonitor()

	// wait for mq runner to stop
	logger.Logmsg(fmt.Sprintf("startup time: %v", time.Since(startup)))
	mqrunner.WaitForRunnerStop()

	logger.Logmsg("exiting")
}
