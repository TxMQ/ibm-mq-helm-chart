package qmgr

import (
	"fmt"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
	"time"
)

func StartMonitor(qmgr string) {
	logger.Logmsg("start-monitor: starting qmgr state monitor")

	go runmonitor(qmgr)
}

func runmonitor(qmgr string) {
	currstatus, err := util.QmgrStatus(qmgr, true)
	if err != nil {
		// alert?
	}

	for {
		time.Sleep(5 * time.Second)

		status, err := util.QmgrStatus(qmgr, true)
		if err != nil {
			// alert?

		} else if status != currstatus {
			logger.Logmsg(fmt.Sprintf("qmgr '%s' status changed from '%s' to '%s'", qmgr, currstatus, status))
			currstatus = status

			// display current status
			_, _ = util.QmgrStatus(qmgr, false)
		}
	}
}