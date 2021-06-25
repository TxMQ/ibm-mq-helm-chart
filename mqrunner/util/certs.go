package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const keyDatabaseStem = "key"
const keydbSuffix = ".kdb"
const rdbSuffix = ".rdb"
const sthSuffix = ".sth"

const _keyfile = "tls.key"
const _certfile = "tls.crt"
const _cafile = "ca.crt"

const _certlabel = "ibmwebspheremq"

//
// ImportCertificates from certdir into keyrepo in ssldir
//
func ImportCertificates(qmgr string) error {

	//
	// /etc/mqm/pki/certs - pki keys and certs
	// /etc/mqm/pki/trust - pki trust roots
	// /etc/mqm/ssl - qmgr key repo directory
	//

	certdir := "/etc/mqm/pki/certs"
	ssldir := "/etc/mqm/ssl"
	trustdir := "/etc/mqm/pki/trust"

	// certs are mounted into the container as secrets
	// with keys tls.key, tls.crt, and ca.crt
	// tls.crt certificate contains cert chain not including root ca

	// create self-signed key pair
	// openssl req -newkey rsa:2048 -nodes -keyout tls.key -subj "/CN=localhost" -x509 -days 3650 -out tls.crt

	// re-create cms keystore
	kdbpath, err := RecreateCMSKeyStore(ssldir)
	if err != nil {
		return err
	}

	// import ca chains into the keystore. ca-chains include self-signed certs.
	err = ImportTrustChains(kdbpath, certdir, trustdir)
	if err != nil {
		return err
	}

	// format cert label
	certlabel := formatCertLabel(qmgr)

	// convert pem key and cert files into p12 format
	p12path, err := PemToP12(certdir, ssldir, certlabel)
	if err != nil {
		return err
	}

	// import p12 file into the keystore
	err = ImportP12(p12path, kdbpath, certlabel)
	if err != nil {
		return err
	}

	// delete p12 file
	err = deleteFile(p12path)
	if err != nil {
		return err
	}

	return nil
}

func RecreateCMSKeyStore(ssldir string) (string, error) {
	return CreateCMSKeyStore(ssldir, true)
}

func expandKeyDatabaseStemName(stem string) (string, string, string) {
	return stem + keydbSuffix, stem + rdbSuffix, stem + sthSuffix
}

func CreateCMSKeyStore(ssldir string, deleteExistingKeystore bool) (string, error) {

	// check if ssldir exists
	_, err := os.Stat(ssldir)
	if err != nil {
		return "", err
	}

	// keyr is key repo file w/o extension
	// keyr expands into 3 files: keyr.kdb, keyr.rdb, keyr.sth
	keydb, rdb, sth := expandKeyDatabaseStemName(keyDatabaseStem)

	keydbpath := filepath.Join(ssldir, keydb)
	rdbpath := filepath.Join(ssldir, rdb)
	sthpath := filepath.Join(ssldir, sth)

	if deleteExistingKeystore {
		for _, fpath := range []string {keydbpath, rdbpath, sthpath} {
			err = os.Remove(fpath)
			if err != nil && !os.IsNotExist(err) {
				return "", err
			}
		}
	}

	// runmqckm -keydb -create -db zorro.kdb -pw password -type cms -stash
	// -rw-------. 1 1000680000 root  88 Jun 23 17:28 key.kdb
	//-rw-------. 1 1000680000 root  80 Jun 23 17:28 key.rdb
	//-rw-------. 1 1000680000 root 193 Jun 23 17:28 key.sth

	// generate password
	password, err := exec.Command("openssl", "rand", "-base64", "14").CombinedOutput()
	if err != nil {
		return "", err
	}

	// create keystore
	err = exec.Command("/opt/mqm/bin/runmqckm", "-keydb", "-create", "-db", keydbpath,
		"-pw", string(password), "-type", "cms", "-stash").Run()

	if err != nil {
		return "", err
	}

	// change access mode for the keyrepo
	// chmod g+rw zorro.*
	// -rw-rw----. 1 1000680000 root  88 Jun 23 17:28 key.kdb
	//-rw-rw----. 1 1000680000 root  80 Jun 23 17:28 key.rdb
	//-rw-rw----. 1 1000680000 root 193 Jun 23 17:28 key.sth

	for _, fpath := range []string {keydbpath, rdbpath, sthpath} {
		err = os.Chmod(fpath, 0660)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}

	return keydbpath, nil
}

func IsSelfSigned(certpath string) (bool, error) {

	out, err := exec.Command("openssl", "x509", "-text", "-in", certpath, "-noout").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("%v\n", string(out))
	}

	cout := string(out)

	issuerIdx := strings.Index(cout, "Issuer:")
	if issuerIdx < 0 { return false, nil }

	colidx := strings.Index(cout[issuerIdx:], ":")
	nlidx := strings.Index(cout[issuerIdx:], "\n")
	issuer := strings.TrimSpace(cout[issuerIdx+colidx+1: issuerIdx+nlidx])

	subjectIdx := strings.Index(cout, "Subject:")
	if subjectIdx < 0 { return false, nil }

	colidx = strings.Index(cout[subjectIdx:], ":")
	nlidx = strings.Index(cout[subjectIdx:], "\n")
	subject := strings.TrimSpace(cout[subjectIdx+colidx+1: subjectIdx+nlidx])

	return issuer == subject, nil
}

func ImportTrustChains(keydbpath, certdir, trustdir string) error {

	// in the certdir we expect:
	// tls.key, tls.crt, [ca.crt]

	// if certdir is mounted from the certificate-manger secret then
	// tls.cert may contain cert chain, terminating before root ca.
	// root ca will be in the ca.crt file.
	// trustdir may conain additional trust chains

	// check if cert directory exists
	_, err := os.Stat(certdir)
	if err != nil {
		return err
	}

	// check if tls.crt exists
	certpath := filepath.Join(certdir, _certfile)

	_, err = os.Stat(certpath)
	if err != nil {
		return err
	}

	selfsigned, err := IsSelfSigned(certpath)
	if err != nil {
		return err
	}

	if selfsigned {

		// add self-signed certificate:
		// runmqckm -cert -add -db filename -stashed -label label -file filename -format ascii

		label := "ssca"
		err = exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-add", "-db", keydbpath, "-stashed",
			"-label", label, "-file", certpath, "-format", "ascii").Run()

		if err != nil {
			return err
		}
	}

	//
	// certdir 'may' contain ca.crt
	//
	capath := filepath.Join(certdir, _cafile)

	_, err = os.Stat(capath)
	if err == nil {

		// add ca certificate:
		// runmqckm -cert -add -db filename -stashed -label label -file filename -format ascii

		label := "ca1"
		err = exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-add", "-db", keydbpath, "-stashed",
			"-label", label, "-file", capath, "-format", "ascii").Run()

	} else if !os.IsNotExist(err) {
		return err
	}

	// todo
	// trust directory may contain trust chains
	// chain: root->ca1->ca2->...->ca

	return nil
}

func PemToP12(certdir, ssldir, certlabel string) (string, error) {

	keypath := filepath.Join(certdir, _keyfile)
	_, err := os.Stat(keypath)
	if err != nil {
		return "", err
	}

	certpath := filepath.Join(certdir, _certfile)
	_, err = os.Stat(certpath)
	if err != nil {
		return "", err
	}

	p12path :=  filepath.Join(ssldir, "qm.p12")
	_, err = os.Stat(ssldir)
	if err != nil {
		return "", err
	}

	// openssl pkcs12 -export -name "label" -out qm.p12 -inkey keyfile -in certfile [-certfile chainfile]
	// -keypbe NONE -certpbe NONE -nomaciter -passout pass:
	// note that we do not include cert chain because cert chains are imported separately

	cmd := exec.Command("/usr/bin/openssl", "pkcs12", "-export", "-name", certlabel, "-out", p12path,
		"-inkey", keypath, "-in", certpath, "-keypbe", "NONE", "-certpbe", "NONE", "-nomaciter",
		"-passout", "pass:")

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	// chanage p12 mode
	err = os.Chmod(p12path, 0660)
	if err != nil {
		return "", err
	}

	return p12path, nil
}

func ImportP12(p12path, kdbpath, certlabel string) error {

	_, err := os.Stat(p12path)
	if err != nil {
		return err
	}

	_, err = os.Stat(kdbpath)
	if err != nil {
		return err
	}

	// runmqckm -cert -import -file ./qm.p12 -pw "" -type pkcs12 -target ./zorro.kdb -target_pw password
	// -target_type cms -label label -new_label qm

	cmd := exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-import", "-file", p12path,
		"-pw", "", "-type", "pkcs12", "-target", kdbpath, "-target_stashed",
		"-target_type", "cms", "-label", certlabel, "-new_label", certlabel)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func formatCertLabel(qmgr string) string {
	return _certlabel + qmgr
}

func deleteFile(file string) error {
	err := os.Remove(file)

	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
