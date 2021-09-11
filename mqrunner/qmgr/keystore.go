package qmgr

import (
	"fmt"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
)

func ImportQmgrKeystore(qmgr string) (bool, string, error) {

	if util.IsEnableTls() {
		// import certs into the keystore
		if util.GetDebugFlag() {
			logger.Logmsg(fmt.Sprintf("importing certificates into qmgr '%s' keystore", qmgr))
		}

		if keypath, err := util.ImportCertificates(qmgr); err != nil {
			// 2021/09/11 00:15:49 [qmgr.ImportQmgrKeystore:18] wait: no child processes
			logger.Logmsg(err)
			return true, "", fmt.Errorf("error importing certificates, %v", err)

		} else {
			// tls enabled and keypath
			return  true, keypath, nil
		}

	} else {
		if util.GetDebugFlag() {
			logger.Logmsg("TLS is not enabled")
		}
		// tls is not enabled
		return false, "", nil
	}
}
