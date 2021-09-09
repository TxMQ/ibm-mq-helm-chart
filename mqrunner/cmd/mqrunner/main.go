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

	// startup role: 0/1
	// from the hostname (k8s) or command line arg (docker)

	logger.Runlogger()
	probe.StartProbes(qmname)

	if qmgr.IsStartupLeader() {
		// new package
		if err := qmgr.CreateDirectories(); err != nil {
			logger.Logmsg(err)
		}

		if err := util.MergeMqscFiles(); err != nil {
			logger.Logmsg(err)
		}

		if err := util.MergeGitConfigFiles(); err != nil {
			logger.Logmsg(err)
		}

		mqrunner.StartMqrunner()
		mqrunner.WaitForRunnerReady()

		if err := qmgr.CreateQmgr(qmname); err != nil {
			logger.Logmsg(err)
		}

		// tail logs
		qmgr.TailLogs(qmname)

		if err := qmgr.ImportQmgrKeystore(qmname); err != nil {
			logger.Logmsg(err)
		}

		if err := qmgr.StartQmgr(qmname); err != nil {
			logger.Logmsg(err)
		}

		if err := util.ApplyStartupConfig(qmname); err != nil {
			logger.Logmsg(err)
		}

	} else {
		mqrunner.StartMqrunner()
		mqrunner.WaitForRunnerReady()

		// wait for qmgr to be created by the leader
		_ = qmgr.WaitForQmgrCreate(qmname)

		// tail logs
		qmgr.TailLogs(qmname)

		if err := qmgr.StartQmgr(qmname); err != nil {
			logger.Logmsg(err)
		}
	}

	// get running role (active, standby)

	// start perf-monitor
	perfmon.StartPerfMonitor()

	// start webc
	mqwebc.StartWebconsole()

	// wait for mq runner to stop
	logger.Logmsg(fmt.Sprintf("startup time: %v", time.Since(startup)))
	mqrunner.WaitForRunnerStop()

	logger.Logmsg("exiting")
}
