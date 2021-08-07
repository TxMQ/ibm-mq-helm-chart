package util

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestCreateQmgr(t *testing.T) {
	qmgr := "qm9"

	err := CreateQmgr(qmgr)
	if err != nil {
		fmt.Printf("create-qmgr: %v\n", err)
	} else {
		fmt.Printf("queue manager %s created\n", qmgr)
	}
}

func TestStartStopQmgr(t *testing.T) {
	qmgr := "qm9"

	err := StartQmgr(qmgr)
	if err != nil {
		fmt.Printf("start-qmgr: %v\n", err)
	} else {
		fmt.Printf("queue manager %s started\n", qmgr)
	}

	err = StopQmgr(qmgr)
	if err != nil {
		fmt.Printf("stop-qmgr: %v\n", err)
	} else {
		fmt.Printf("queue manager %s stopped\n", qmgr)
	}
}

func TestQmgrExists(t *testing.T) {

	exists, err := QmgrExists("qm2")
	if err != nil {
		println("error...")
	}

	if exists {
		println("exists...")
	} else {
		println("does not exist")
	}

	exists, err = QmgrExists("qm3")
	if err != nil {
		println("error...")
	}

	if exists {
		println("exists...")
	} else {
		println("does not exist")
	}
}

func TestIsQmgrRunning(t *testing.T) {
	running, err := IsQmgrRunning("qm2")
	if err != nil {
		println("error")
	} else if running {
		println("running...")
	} else {
		println("not running...")
	}

	running, err = IsQmgrRunning("qm3")
	if err != nil {
		println("error")
	} else if running {
		println("running...")
	} else {
		println("not running...")
	}
}

func TestGetCertLabel(t *testing.T) {

	label, err := GetCertLabel("qm2")

	if err != nil {
		println("test error...")
	} else {
		println(label)
	}
}

func TestSetSslKeyRepo(t *testing.T) {
	keyr, err := GetSslKeyRepo("qm2")

	if err != nil {
		println("test error...")
	} else {
		println(keyr)
	}

	keyr, err = GetSslKeyRepo("qm3")

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		println(keyr)
	}
}

func TestCreateCMSKeyStore(t *testing.T) {

	u, _ := user.Current()
	ssldir := filepath.Join(u.HomeDir, "etc/mqm/ssl")

	path, err := CreateCMSKeyStore(ssldir, true)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("key store path: %s\n", path)
	}
}

func TestIsSelfSigned(t *testing.T) {

	u, _ := user.Current()
	certpath := filepath.Join(u.HomeDir, "etc/mqm/pki/tls.crt")

	_, _, ss, err := IsSelfSigned(certpath)

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("self-signed: %v\n", ss)
	}
}

func TestImportTrustChains(t *testing.T) {
	u, _ := user.Current()
	ssldir := filepath.Join(u.HomeDir, "etc/mqm/ssl")

	keydbpath, err := CreateCMSKeyStore(ssldir, true)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	certdir := filepath.Join(u.HomeDir, "etc/mqm/pki/cert")
	trustdir := filepath.Join(u.HomeDir, "etc/mqm/pki/trust")

	// TODO: fix test case
	err = ImportTrustCerts(keydbpath, certdir, "", trustdir)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestPemToP12(t *testing.T) {
	u, _ := user.Current()
	ssldir := filepath.Join(u.HomeDir, "etc/mqm/ssl")
	certdir := filepath.Join(u.HomeDir, "etc/mqm/pki/cert")
	certlabel := "ibmwebspheremqqm2"

	// TODO: fix test case
	p12file, err := PemToP12("", certdir, ssldir, certlabel)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", p12file)
	}
}

func TestImportP12(t *testing.T) {
	u, _ := user.Current()

	p12path := filepath.Join(u.HomeDir, "etc/mqm/ssl/qm.p12")
	kdbpath := filepath.Join(u.HomeDir, "etc/mqm/ssl/key.kdb")
	certlabel := "ibmwebspheremqqm2"

	err := ImportP12(p12path, kdbpath, certlabel)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", "imported")
	}
}

func TestCreateDirectories(t *testing.T) {

	err := CreateDirectories()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("directories created...")
	}
}
