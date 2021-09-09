package mqwebc

import (
	"os"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
	"szesto.com/mqrunner/webmq"
)

func isStartMqweb() bool {
	return os.Getenv("MQ_START_MQWEB") == "1"
}

func isConfigureMqweb() bool {
	return os.Getenv("MQ_CONFIGURE_MQWEB") == "1"
}

func StartWebconsole() {
	// configure webconsole
	if isStartMqweb() || isConfigureMqweb() {
		go startwebc()

	} else {
		logger.Logmsg("webconsole is off")
	}
}

func startwebc() {
	logger.Logmsg("configuring webconsole")

	err := webmq.ConfigureWebconsole()
	if err != nil {
		// log error
		logger.Logmsg(err)

		if isStartMqweb() {
			logger.Logmsg("web console configuration failed, web console will not be started")
		}

	} else {
		if isStartMqweb() {
			logger.Logmsg("starting mq web console")
			logger.Logmsg("web console will connect to ldap, if taking too long check ldap server")

			err = util.StartMqweb()
			if err != nil {
				// log error
				logger.Logmsg(err)
			}
		}
	}
}