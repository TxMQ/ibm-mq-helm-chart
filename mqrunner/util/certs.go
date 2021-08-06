package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const _keyDatabaseStem = "key"
const _keydbSuffix = ".kdb"
const _rdbSuffix = ".rdb"
const _sthSuffix = ".sth"

const _keyfile = "tls.key"
const _certfile = "tls.crt"
const _cafile = "ca.crt"

const _certlabel = "ibmwebspheremq"

const _ssldir = "/etc/mqm/ssl"
const _certdir = "/etc/mqm/pki/cert"
const _trustdir = "/etc/mqm/pki/trust"

func IsEnableTls() bool {

	if vaultTls := os.Getenv("VAULT_ENABLE_TLS"); len(vaultTls) > 0 {
		if vaultTls == "true" || vaultTls == "1" {
			return true
		}
	}

	if mqTls := os.Getenv("MQ_ENABLE_TLS_NO_VAULT"); len(mqTls) > 0 {
		if mqTls == "true" || mqTls == "1" {
			return true
		}
	}

	return false
}

func SetQmgrKeyRepoLocation(qmgr string) error {
	sslkeyr := filepath.Join(GetSsldir(qmgr), _keyDatabaseStem)
	return SetSslKeyRepo(qmgr, sslkeyr)
}

func getTrustdir() string {
	return _trustdir
}

func GetSsldir(qmgr string) string {
	if len(qmgr) > 0 {
		// queue manager default
		return fmt.Sprintf("/var/mqm/qmgrs/%s/ssl", qmgr)

	} else {
		// web console
		return _ssldir
	}
}

func GetTlsKeyPath() string {
	if keypath := os.Getenv("VAULT_TLS_KEY_INJECT_PATH"); len(keypath) > 0 {
		return keypath
	} else {
		return filepath.Join(_certdir, _keyfile)
	}
}

func GetTlsCertPath() string {
	if certpath := os.Getenv("VAULT_TLS_CERT_INJECT_PATH"); len(certpath) > 0 {
		return certpath
	} else {
		return filepath.Join(_certdir, _certfile)
	}
}

func GetTlsCaPath() string {
	if capath := os.Getenv("VAULT_TLS_CA_INJECT_PATH"); len(capath) > 0 {
		return capath
	} else {
		return filepath.Join(_certdir, _cafile)
	}
}

//
// ImportCertificates from certdir into keyrepo in ssldir
//
func ImportCertificates(qmgr string) error {

	// cert dir with key, cert, ca cert
	//certdir := getCertdir()

	keypath := GetTlsKeyPath()
	certpath := GetTlsCertPath()
	capath := GetTlsCaPath()

	// trusted certs
	trustdir := getTrustdir()

	// key store directory
	ssldir := GetSsldir(qmgr)

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

	// import ca certs into the keystore. ca-certs include self-signed certs.
	err = ImportTrustCerts(kdbpath, certpath, capath, trustdir)
	if err != nil {
		return err
	}

	// format cert label
	certlabel := formatCertLabel(qmgr)

	// convert pem key and cert files into p12 format
	p12path, err := PemToP12(keypath, certpath, ssldir, certlabel)
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
	return stem + _keydbSuffix, stem + _rdbSuffix, stem + _sthSuffix
}

func CreateCMSKeyStore(ssldir string, deleteExistingKeystore bool) (string, error) {

	// check if ssldir exists
	_, err := os.Stat(ssldir)
	if err != nil {
		return "", err
	}

	// keyr is key repo file w/o extension
	// keyr expands into 3 files: keyr.kdb, keyr.rdb, keyr.sth
	keydb, rdb, sth := expandKeyDatabaseStemName(_keyDatabaseStem)

	keydbpath := filepath.Join(ssldir, keydb)
	rdbpath := filepath.Join(ssldir, rdb)
	sthpath := filepath.Join(ssldir, sth)

	log.Printf("create-cms-keystore-1: keydbpath = %s, rdbpath = %s, sthpath = %s\n", keydbpath, rdbpath, sthpath)

	log.Printf("%s\n", "create-cms-keystore-2: deleting existing keystore")

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

	log.Printf("%s\n", "create-cms-keystore-3: generating keystore password")

	// generate password
	password, err := exec.Command("openssl", "rand", "-base64", "14").CombinedOutput()
	if err != nil {
		return "", err
	}

	log.Printf("create-cms-keystore-4: creating keystore %s\n", keydbpath)

	// create keystore
	out, err := exec.Command("/opt/mqm/bin/runmqckm", "-keydb", "-create", "-db", keydbpath,
		"-pw", string(password), "-type", "cms", "-stash").CombinedOutput()

	if err != nil {
		if out != nil {
			return "", fmt.Errorf("%s\n", string(out))
		} else {
			return "", err
		}
	}

	// change access mode for the keyrepo
	// chmod g+rw zorro.*
	// -rw-rw----. 1 1000680000 root  88 Jun 23 17:28 key.kdb
	//-rw-rw----. 1 1000680000 root  80 Jun 23 17:28 key.rdb
	//-rw-rw----. 1 1000680000 root 193 Jun 23 17:28 key.sth

	log.Printf("create-cms-keystore-5: changing keystore permissions\n")

	for _, fpath := range []string {keydbpath, rdbpath, sthpath} {
		err = os.Chmod(fpath, 0660)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}

	log.Printf("create-cms-keystore-6: keystore %s created\n", keydbpath)

	return keydbpath, nil
}

func IsSelfSigned(certpath string) (string, string, bool, error) {

	out, err := exec.Command("openssl", "x509", "-text", "-in", certpath, "-noout").CombinedOutput()
	if err != nil {
		return "", "", false, fmt.Errorf("%v\n", string(out))
	}

	cout := string(out)

	issuerIdx := strings.Index(cout, "Issuer:")
	if issuerIdx < 0 { return "", "", false, nil }

	colidx := strings.Index(cout[issuerIdx:], ":")
	nlidx := strings.Index(cout[issuerIdx:], "\n")
	issuer := strings.TrimSpace(cout[issuerIdx+colidx+1: issuerIdx+nlidx])

	subjectIdx := strings.Index(cout, "Subject:")
	if subjectIdx < 0 { return "", "", false, nil }

	colidx = strings.Index(cout[subjectIdx:], ":")
	nlidx = strings.Index(cout[subjectIdx:], "\n")
	subject := strings.TrimSpace(cout[subjectIdx+colidx+1: subjectIdx+nlidx])

	return subject, issuer,  issuer == subject, nil
}

func ImportTrustCerts(keydbpath, certpath, capath, trustdir string) error {

	// check if cert path exists
	log.Printf("import-trust-chains-1: expecting certificate %s\n", certpath)

	_, err := os.Stat(certpath)
	if err != nil {
		return err
	}

	subject, issuer, selfsigned, err := IsSelfSigned(certpath)
	if err != nil {
		return err
	}

	log.Printf("import-trust-chains-2: found certificate %s. subject='%s', issuer='%s'\n",
		certpath, subject, issuer)

	if selfsigned {

		log.Printf("import-trust-chains-3: certificate %s is self-signed, importing into key db %s\n",
			certpath, keydbpath)

		// add self-signed certificate:

		label := "ssca"

		if err = addcert(keydbpath, label, certpath); err != nil {
			return err
		}
	}

	//
	// certdir 'may' contain ca.crt
	//
	_, err = os.Stat(capath)
	if err == nil {

		// add ca certificate:
		// runmqckm -cert -add -db filename -stashed -label label -file filename -format ascii

		label := "ca"

		log.Printf("import-trust-chains-4: importing ca cert %s into key db %s, label %s\n",
			capath, keydbpath, label)

		if err = addcert(keydbpath, label, capath); err != nil {
			return nil
		}

	} else if !os.IsNotExist(err) {
		return err
	}

	// trust directory may contain trust certs, *.crt
	_, err = os.Stat(trustdir)
	if err == nil {
		trustpems, err := ReadDir(trustdir, "crt")
		if err != nil {
			return err
		}

		if len(trustpems) > 0 {
			for idx, trustpem := range trustpems {

				label := fmt.Sprintf("trust%d", idx)

				log.Printf("import-trust-chains-5: importing ca cert %s into key db %s, label %s\n",
					trustpem, keydbpath, label)

				if err = addcert(keydbpath, label, trustpem); err != nil {
					return err
				}
			}
		}

	} else if !os.IsNotExist(err) {
		return nil
	}

	return nil
}

func addcert(keydbpath, label, certpath string) error {

	// runmqckm -cert -add -db filename -stashed -label label -file filename -format ascii

	out, err := exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-add", "-db", keydbpath, "-stashed",
		"-label", label, "-file", certpath, "-format", "ascii").CombinedOutput()

	if err != nil {
		if out != nil {
			return fmt.Errorf("%s\n", string(out))
		} else {
			return err
		}
	}

	return nil
}

func PemToP12(keypath, certpath, ssldir, certlabel string) (string, error) {

	//keypath := filepath.Join(certdir, _keyfile)
	_, err := os.Stat(keypath)
	if err != nil {
		return "", err
	}

	//certpath := filepath.Join(certdir, _certfile)
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

	log.Printf("pem-to-p12-1: converting key %s and cert %s into p12 %s\n", keypath, certpath, p12path)

	out, err := exec.Command("/usr/bin/openssl", "pkcs12", "-export", "-name", certlabel, "-out", p12path,
		"-inkey", keypath, "-in", certpath, "-keypbe", "NONE", "-certpbe", "NONE", "-nomaciter",
		"-passout", "pass:").CombinedOutput()

	if err != nil {
		if out != nil {
			return "", fmt.Errorf("%s\n", string(out))
		} else {
			return "", err
		}
	}

	log.Printf("pem-to-p12-2: changing %s permissions", p12path)

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

	log.Printf("import-p12-1: importing p12 file %s into key db %s with cert label '%s'\n",
		p12path, kdbpath, certlabel)

	out, err := exec.Command("/opt/mqm/bin/runmqckm", "-cert", "-import", "-file", p12path,
		"-pw", "", "-type", "pkcs12", "-target", kdbpath, "-target_stashed",
		"-target_type", "cms", "-label", certlabel, "-new_label", certlabel).CombinedOutput()

	if err != nil {
		if out != nil {
			return fmt.Errorf("%s\n", string(out))
		} else {
			return err
		}
	}

	log.Printf("import-p12-2: imported p12 file %s into key db %s with cert label '%s'\n",
		p12path, kdbpath, certlabel)

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
