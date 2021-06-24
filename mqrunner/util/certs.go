package util

import (
	"os"
	"os/exec"
)

type TlsInfo struct {
	certdir string // pem secrets mounted on this directory
	keyfile string // key file in cert directory: tls.key
	certfile string // cert file in cert directory: tls.crt
	cafile string // ca file in cert directory
	ssldir string // queue manager ssl directory
	keyrepo string // key repo in qm ssl directory
	p12file string // p12 file in ssl directory
}

//
// ImportCertificates from certdir into keyrepo in ssldir
//
func ImportCertificates(qmgr string) error {

	//
	// /etc/mqm/pki/certs - pki keys and certs
	// /etc/mqm/pki/trust - pki trust roots
	// /etc/mqm/ssl - util key repo directory
	//

	certdir := "/etc/mqm/pki/certs"
	ssldir := "/etc/mqm/ssl"
	keyrepo := "key.kdb"

	// certs are mounted into the container as secrets
	// with keys tls.key, tls.crt, and ca.crt
	// tls.crt certificate contains cert chain not including root ca

	// create self-signed key pair
	// openssl req -newkey rsa:2048 -nodes -keyout tls.key -subj "/CN=localhost" -x509 -days 3650 -out tls.crt

	// re-create cms keystore
	keyrepoPath, err := RecreateCMSKeyStore(ssldir, keyrepo)
	if err != nil {
		return err
	}

	// import ca chains into the keystore. ca-chains include self-signed certs.
	err = importCAChain(keyrepoPath, certdir)
	if err != nil {
		return err
	}

	// call a function
	certlabel := "imbwebshperemq" + qmgr

	// convert pem key and cert files into p12 format
	p12file, err := pemToP12(certdir, ssldir, certlabel)
	if err != nil {
		return err
	}

	// import p12 file into the keystore
	err = importP12(p12file, keyrepoPath, certlabel)
	if err != nil {
		return err
	}

	return nil
}

func RecreateCMSKeyStore(ssldir, keyrepo string) (string, error) {
	return CreateCMSKeyStore(ssldir, keyrepo, true)
}

func CreateCMSKeyStore(ssldir, keyrepo string, deleteOldKeystore bool) (string, error) {

	// check if ssldir exists

	// build key repo path
	keyrepoPath := ssldir + string(os.PathSeparator) + keyrepo

	// delete old keystore if requested
	if deleteOldKeystore == true {

	}

	// runmqckm -keydb -create -db zorro.kdb -pw password -type cms -stash
	// -rw-------. 1 1000680000 root  88 Jun 23 17:28 zorro.kdb
	//-rw-------. 1 1000680000 root  80 Jun 23 17:28 zorro.rdb
	//-rw-------. 1 1000680000 root 193 Jun 23 17:28 zorro.sth

	// todo: generate password

	cmd := exec.Command("/opt/mqm/bin/runmqckm", "-keydb", "-create", "-db", keyrepoPath,
		"-pw", "password", "-type", "cms", "-stash")

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// change access mode for the keyrepo
	// chmod g+rw zorro.*
	// -rw-rw----. 1 1000680000 root  88 Jun 23 17:28 zorro.kdb
	//-rw-rw----. 1 1000680000 root  80 Jun 23 17:28 zorro.rdb
	//-rw-rw----. 1 1000680000 root 193 Jun 23 17:28 zorro.sth

	return keyrepoPath, nil
}

func importCAChain(keyrepoPath, certdir string) error {

	// certdir contains files: tls.key, tls.crt, ca.crt

	// todo: check if tls.crt is self-signed cert
	// for now assume it is

	// self signed cert must be imported with this function
	// runmqckm -cert -add -db ./zorro.kdb -stashed -label zorro -file ./tls.crt -format ascii

	label := "ca"
	cert := "tls.crt"
	certpath := certdir + string(os.PathSeparator) + cert

	cmd := exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-add", "-db", keyrepoPath, "-stashed",
		"-label", label, "-file", certpath, "-format", "ascii")

	err := cmd.Run()
	if err != nil {
		return err
	}

	// chain: root->ca1->ca2->...->ca
	// root->ca1->ca

	// to import:
	// runmqckm -cert -add -db filename -stashed -label label
	// -file filename -format ascii

	// runmqakm -cert -add -db filename -stashed -label label
	// -file filename -format ascii -fips

	return nil
}

func pemToP12(certdir, ssldir, certlabel string) (string, error) {

	keyfile := "tls.key"
	keypath := certdir + string(os.PathSeparator) + keyfile

	certfile := "tls.crt"
	certpath := certdir + string(os.PathSeparator) + certfile

	//cafile := "ca.crt"
	//capath := certdir + string(os.PathSeparator) + cafile

	p12path := ssldir + string(os.PathSeparator) + "qm.p12"

	// openssl pkcs12 -export -name "label" -out qm.p12 -inkey keyfile -in certfile -certfile chainfile
	// -keypbe NONE -certpbe NONE -nomaciter -passout pass:

	cmd := exec.Command("/usr/bin/openssl", "pkcs12", "-export", "-name", certlabel, "-out", p12path,
		"-inkey", keypath, "-in", certpath, "-keypbe", "NONE", "-certpbe", "NONE", "-nomaciter",
		"-passout", "pass:")

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return p12path, nil
}

func importP12(p12path, kdbpath, certlabel string) error {

	// runmqckm -cert -import -file ./qm.p12 -pw "" -type pkcs12 -target ./zorro.kdb -target_pw password
	// -target_type cms -label label -new_label qm

	cmd := exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-import", "-file", p12path,
		"-pw", "", "-type", "pkcs12", "-target", kdbpath, "-target_stashed",
		"-target_type", "cms", "-label", certlabel, "-new_label", certlabel)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}