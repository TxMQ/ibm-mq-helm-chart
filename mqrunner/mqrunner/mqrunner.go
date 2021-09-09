package mqrunner

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"syscall"
	"szesto.com/mqrunner/logger"
)

var runnerchan chan int

func StartMqrunner() {
	logger.Logmsg("starting mq runner")
	runnerchan = make(chan int)
	go mqrunner()
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

func mqrunner() {

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
			logger.Logmsg(fmt.Sprintf("%s", "shutting down queue manager"))

			logger.Logmsg("mq runner stopped")
			runnerStopped()
			break
		}
	}
}

func WaitForRunnerStop() {
	logger.Logmsg("waiting for mq runner stop")
	<- runnerchan
}
