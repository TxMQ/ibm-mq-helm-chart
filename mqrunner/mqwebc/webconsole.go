package mqwebc

import (
	"fmt"
	"os"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
	"szesto.com/mqrunner/webmq"
	"time"
)

func IsStartMqweb() bool {
	return os.Getenv("MQ_START_MQWEB") == "1"
}

func IsConfigureMqweb() bool {
	return os.Getenv("MQ_CONFIGURE_MQWEB") == "1"
}

func StartWebconsole() {
	// configure webconsole
	if IsStartMqweb() || IsConfigureMqweb() {
		startwebc()

	} else {
		logger.Logmsg("webconsole is off")
	}
}

func startwebc() {
	t := time.Now()

	logger.Logmsg("configuring web console")

	err := webmq.ConfigureWebconsole()
	if err != nil {
		// log error
		logger.Logmsg(err)

		if IsStartMqweb() {
			logger.Logmsg("web console configuration failed, web console will not be started")
		}

	} else {
		if IsStartMqweb() {
			logger.Logmsg("starting mq web console")
			logger.Logmsg("web console will connect to ldap, if taking too long check ldap server")

			err = util.StartMqweb()
			if err != nil {
				// log error
				logger.Logmsg(err)
			}
		}
	}

	// log elapse time
	logger.Logmsg(fmt.Sprintf("time to import web console keystore: %v", time.Since(t)))
}