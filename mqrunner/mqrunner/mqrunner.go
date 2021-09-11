package mqrunner

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"syscall"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
)

var runnerchan chan int

func StartMqrunner(qmgr string) {
	logger.Logmsg(fmt.Sprintf("starting mq runner for qmgr '%s'", qmgr))
	runnerchan = make(chan int)
	go mqrunner(qmgr)
}

func runnerReady() {
	runnerchan <- 0
}

func WaitForRunnerReady() {
	logger.Logmsg("waiting for mq runner ready")
	<- runnerchan
}

func runnerStopped() {
	runnerchan <- 1
}

func mqrunner(qmgr string) {

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigcld := make(chan os.Signal)
	signal.Notify(sigcld, syscall.SIGCHLD)

	logger.Logmsg("mq runner ready")
	runnerReady()

	for {
		select {
		case <- sigcld:
			var waitStatus unix.WaitStatus
			/*pid err*/ _, _ = unix.Wait4(-1, &waitStatus, unix.WNOHANG, nil)

		case <- sigterm:
			logger.Logmsg(fmt.Sprintf("%s", "received sigterm, exiting"))

			// shutdown queue manager
			if running, err := util.IsQmgrRunning(qmgr, false); err == nil && running {
				logger.Logmsg(fmt.Sprintf("shutting down queue manager '%s'", qmgr))
				_ = util.StopQmgr(qmgr)
			}

			logger.Logmsg("mq runner stopped")
			runnerStopped()
			break
		}
	}
}

func WaitForRunnerStop() {
	logger.Logmsg("mq runner forever")
	<- runnerchan
}
