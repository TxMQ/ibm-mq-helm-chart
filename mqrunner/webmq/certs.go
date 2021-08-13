package webmq

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"szesto.com/mqrunner/util"
)

const _keyfile = "tls.key"
const _certfile = "tls.crt"
const _cafile = "ca.crt"

const _certlabel = "default"

func ImportWebconsoleCerts() (string, string, error) {

	keypath := util.GetTlsKeyPath()
	certpath := util.GetTlsCertPath()
	capath := util.GetTlsCaPath()

	ssldir := util.GetSsldir("")

	return RecreateWebmqKeystore(ssldir, keypath, certpath, capath)
}

func RecreateWebmqKeystore(ssldir, keypath, certpath, capath string) (string, string, error) {
	return CreateWebmqKeystore(ssldir, keypath, certpath, capath, true)
}

func CreateWebmqKeystore(ssldir, keypath, certpath, capath string, deleteExistingKeystore bool) (string, string, error) {

	debug := util.GetDebugFlag()

	// check if ssldir exists
	_, err := os.Stat(ssldir)
	if err != nil {
		return "", "", err
	}

	// we expect to find key, cert and [ca] files

	_, err = os.Stat(keypath)
	if err != nil {
		return "", "", err
	}

	_, err = os.Stat(certpath)
	if err != nil {
		return "", "", err
	}

	iscapath := true

	_, err = os.Stat(capath)
	if err != nil {
		iscapath = false
	}

	p12path := filepath.Join(ssldir, "webmq.p12")

	if debug {
		log.Printf("create-web-keystore-1: re-creating keystore %s\n", p12path)
	}

	if deleteExistingKeystore {
		if debug {
			log.Printf("create-web-keystore-2: deleting keystore %s\n", p12path)
		}

		_, err := os.Stat(p12path)
		if err != nil && os.IsNotExist(err) {
			log.Printf("create-web-keystore-20: keystore %s does not exist\n", p12path)
		} else if err != nil {
			return "", "", err
		} else {
			if debug {
				log.Printf("create-web-keystore-3: deleting keystore '%s'\n", p12path)
				log.Printf("create-web-keystore-4: rm -f %s\n", p12path)
			}

			out, err := exec.Command("rm", "-f", p12path).CombinedOutput()
			if err != nil {
				if out != nil {
					return "", "", fmt.Errorf("create-web-keystore: %s\n", out)
				} else {
					return "", "", err
				}
			}

			if out != nil {
				log.Printf("create-web-keystore-5: %s\n", out)
			}
		}
	}

	// use openssl to create p12 file from key,cert,chain input
	// p12 file will contain private key, cert and complete cert chain including root.

	// generate random password
	// openssl rand -base64 14 > keystore.password
	if debug {
		log.Printf("create-web-keystore-6: %s\n", "openssl rand -base64 14")
	}

	passbytes, err := exec.Command("openssl", "rand", "-base64", "14").CombinedOutput()
	if err != nil {
		return "","", err
	}

	// strip newline character at the end of the password
	password := strings.TrimSuffix(string(passbytes), "\n")

	if debug {
		log.Printf("create-web-keystore-7: %s\n", "/opt/mqm/web/bin/securityUtility encode password")
	}

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

		if debug {
			log.Printf("create-web-keystore-8: /usr/bin/openssl pkcs12 -export -name %s -out %s -inkey %s -in %s -certfile %s -password pass:password\n",
				_certlabel, p12path, keypath, certpath, capath)
		}

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

		if debug {
			log.Printf("create-web-keystore-9: /usr/bin/openssl pkcs12 -export -name %s -out %s -inkey %s -in %s -password pass:password\n",
				_certlabel, p12path, keypath, certpath)
		}

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
	if debug {
		log.Printf("create-web-keystore-10: chmod 0666 %s\n", p12path)
	}

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
