package qmgr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
	"time"
)

func IsStartupLeader() bool {
	// 0 - leader || !0 && !1 - single instance
	if util.IsMultiInstance1() {
		return true

	} else if util.IsMultiInstance2() {
		return false

	} else {
		return true
	}
}

func StartupRole() string {
	if util.IsMultiInstance1() {
		return "multi-instance active"

	} else if util.IsMultiInstance2() {
		return "multi-instance standby"

	} else {
		return "standalone active"
	}
}

func IsRunningRoleActive(qmgr string) bool {
	status, err := util.QmgrStatus(qmgr, false)
	if err != nil {
		logger.Logmsg(err)
	}
	return status == util.QmgrStatusEnumRunning()
}

func IsRunningRoleStandby(qmgr string) bool {
	status, err := util.QmgrStatus(qmgr, false)
	if err != nil {
		logger.Logmsg(err)
	}
	return status == util.QmgrStatusEnumStandby()
}

func IsRunningRoleRunningElsewhere(qmgr string) bool {
	status, err := util.QmgrStatus(qmgr, false)
	if err != nil {
		logger.Logmsg(err)
	}
	return status == util.QmgrStatusEnumElsewhere()
}

func CreateDirectories() error {
	if util.GetDebugFlag() {
		logger.Logmsg("creating directories")
	}

	t := time.Now()
	defer logger.Logmsg(fmt.Sprintf("time to create directories: %v", time.Since(t)))

	out, err := runcmd("/opt/mqm/bin/crtmqdir_setuid", "-f", "-a")
	if err != nil {
		logger.Logmsg(err)

	} else if util.GetDebugFlag() && len(out) > 0 {
		logger.Logmsg(out)
	}

	return err
}

func WaitForQmgrCreate(qmgr string) error {

	logger.Logmsg(fmt.Sprintf("waiting for the leader to create qmgr '%s'", qmgr))

	// look for /var/md/<qm>/qm.ini file
	qmini := fmt.Sprintf("/var/md/%s/qm.ini", qmgr)

	for {
		if _, err := os.Stat(qmini); err != nil {
			// put upper limit on total wait time
			time.Sleep(5 * time.Second)

		} else {
			return nil
		}
	}
}

func CreateQmgr(qmgr string, retryCount int) error {
	logger.Logmsg(fmt.Sprintf("creating queue manager '%s', retry count %d", qmgr, retryCount))

	t := time.Now()
	defer logger.Logmsg(fmt.Sprintf("create queue manager time: %v", time.Since(t)))

	debug := util.GetDebugFlag()

	// check if qmgr already configured
	qmconf, msg, err := util.QmgrConf(qmgr)
	//qmconf, err := util.QmgrExists(qmgr)
	if err != nil {
		// log error message, continue as if qmgr does not exist
		logger.Logmsg(err)

	} else {
		logger.Logmsg(fmt.Sprintf("%s", msg))
	}

	if qmconf == false {
		if debug {
			logger.Logmsg(fmt.Sprintf("creating queue manager '%s'", qmgr))
		}

		// create queue manager
		if err = util.CreateQmgr(qmgr, false); err != nil {
			return err

		} else {
			logger.Logmsg(fmt.Sprintf("queue manager '%s' created", qmgr))
		}

	} else {
		if debug {
			logger.Logmsg(fmt.Sprintf("queue manager '%s' exists", qmgr))
		}
	}

	return nil
}

func StartQmgr(qmgr string) error {
	logger.Logmsg(fmt.Sprintf("starting queue manager %s", qmgr))

	var out string
	var err error

	if util.IsMultiInstance1() || util.IsMultiInstance2() {
		// add -x argument for the multi-instance start
		out, err = runcmd("/opt/mqm/bin/strmqm", "-x", qmgr)

	} else {
		// start queue manager
		out, err = runcmd("/opt/mqm/bin/strmqm", qmgr)
	}

	if err != nil {
		logger.Logmsg(err)
		cerr := strings.TrimSpace(fmt.Sprintf("%v", err))

		// standby instance started, exit status 30
		if strings.HasSuffix(cerr, "exit status 30") {
			return nil

		} else {
			return fmt.Errorf("failed to start qmgr %s", qmgr)
		}

	} else if len(out) > 0 && util.GetDebugFlag() {
		logger.Logmsg(out)
	}

	return nil
}

func TailLogs(qmgr string) {

	// tail system log: /var/mqm/errors/AMQERR01.LOG
	util.TailMqLog()

	// tail qmgr log: /var/mqm/qmgrs/{qmgr}/errors/AMQERR01.LOG
	util.TailQmgrLog(qmgr)
}

func runcmd(cmd string, args ...string) (string, error) {
	if util.GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("%s %s", cmd, strings.Join(args, " ")))
	}

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		if len(string(out)) > 0 {
			cerr := string(out)
			return "", fmt.Errorf("%s%v", cerr, err)
		} else {
			return "", err
		}
	}

	cout := string(out)
	return cout, nil
}