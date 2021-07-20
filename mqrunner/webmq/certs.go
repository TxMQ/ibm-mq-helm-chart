package webmq

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const _keyfile = "tls.key"
const _certfile = "tls.crt"
const _cafile = "ca.crt"

const _certlabel = "default"

func ImportWebconsoleCerts() (string, string, error) {

	// we assume ssl and cert directories
	// if not, we can pass these values as agruments

	const ssldir = "/etc/mqm/ssl"
	const certdir = "/etc/mqm/pki/cert"

	return RecreateWebmqKeystore(ssldir, certdir)
}

func RecreateWebmqKeystore(ssldir, certdir string) (string, string, error) {
	return CreateWebmqKeystore(ssldir, certdir, true)
}

func CreateWebmqKeystore(ssldir, certdir string, deleteExistingKeystore bool) (string, string, error) {

	// check if ssldir exists
	_, err := os.Stat(ssldir)
	if err != nil {
		return "", "", err
	}

	// check if certdir exists
	_, err = os.Stat(ssldir)
	if err != nil {
		return "", "", err
	}

	// we expect to find key, cert and [ca] files in certdir
	keypath := filepath.Join(certdir, _keyfile)

	_, err = os.Stat(keypath)
	if err != nil {
		return "", "", err
	}

	certpath := filepath.Join(certdir, _certfile)

	_, err = os.Stat(certpath)
	if err != nil {
		return "", "", err
	}

	iscapath := true
	capath := filepath.Join(certdir, _cafile)

	_, err = os.Stat(capath)
	if err != nil {
		iscapath = false
	}

	p12path := filepath.Join(ssldir, "webmq.p12")

	if deleteExistingKeystore {
		// delete existing keystore
	}

	// use openssl to create p12 file from key,cert,chain input
	// p12 file will contain private key, cert and complete cert chain including root.

	// generate random password
	// openssl rand -base64 14 > keystore.password
	passbytes, err := exec.Command("openssl", "rand", "-base64", "14").CombinedOutput()
	if err != nil {
		return "","", err
	}

	// strip newline character at the end of the password
	password := strings.TrimSuffix(string(passbytes), "\n")

	// save password into the file?

	// encode password
	encbytes, err := exec.Command("/opt/mqm/web/bin/securityUtility", "encode", string(password)).CombinedOutput()
	if err != nil {
		if encbytes != nil {
			return "","", fmt.Errorf("%s\n", string(encbytes))
		} else {
			return "","", err
		}
	}

	encpass := strings.TrimSuffix(string(encbytes), "\n")

	// we expect to find key, cert and [ca] files in certdir

	// openssl pkcs12 -export -name name -out p12path -inkey tls.key -in tls.crt [-certfile ca.crt] -password pass:password

	if iscapath {

		out, err := exec.Command("/usr/bin/openssl", "pkcs12", "-export", "-name", _certlabel, "-out", p12path,
			"-inkey", keypath, "-in", certpath, "-certfile", capath, "-password", "pass:" + string(password)).CombinedOutput()

		if err != nil {
			if out != nil {
				return "","", fmt.Errorf("%s\n", string(out))
			} else {
				return "","", err
			}
		}

	} else {

		out, err := exec.Command("/usr/bin/openssl", "pkcs12", "-export", "-name", _certlabel, "-out", p12path,
			"-inkey", keypath, "-in", certpath, "-password", "pass:" + string(password)).CombinedOutput()

		if err != nil {
			if out != nil {
				return "","", fmt.Errorf("%s\n", string(out))
			} else {
				return "","", err
			}
		}
	}

	// change p12 file mode
	err = os.Chmod(p12path, 0666)
	if err != nil {
		return "", "", err
	}

	// keytool -list -rfc -keystore p12path -storetype PKCS12 -storepass `cat keystore.password`
	// keytool will show an alias with the name passed to openssl -name clarg

	// use keytool to update p12 keystore with cert chain.
	// keytool -import -trustcacerts -alias alias_to_be_updated -file chain.pem -keystore keystore.p12

	return p12path, encpass, nil
}
