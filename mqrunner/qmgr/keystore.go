package qmgr

import (
	"fmt"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
	"time"
)

func ImportQmgrKeystore(qmgr string) error {
	t := time.Now()
	defer logger.Logmsg(fmt.Sprintf("time to import qmgr keystore: %v", time.Since(t)))

	if util.IsEnableTls() {
		// import certs into the keystore
		if util.GetDebugFlag() {
			logger.Logmsg(fmt.Sprintf("importing certificates into qmgr '%s' keystore", qmgr))
		}

		if err := util.ImportCertificates(qmgr); err != nil {
			logger.Logmsg(err)
			return fmt.Errorf("error importing certificates, %v", err)
		}

	} else {
		if util.GetDebugFlag() {
			logger.Logmsg("TLS is not enabled")
		}
	}

	return nil
}
