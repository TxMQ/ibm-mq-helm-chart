package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"log"
	"strings"
	"time"
)

type Env struct {

	qmgr string // queue manager name
	conname string // host(port)
	channel string // channel name
	cipherspec string // cipher svrconn configuration
	username string // user name
	password string // password
	qname string // queue name

	// The keystore contains at least the certificate to verify the qmgr's cert (usually from
	// a Certificate Authority) and optionally the client's own certificate.
	// We could also optionally specify which certificate represents the client, based on its label
	// but don't need to do this when using the MQSCA_OPTIONAL flag.
	//

	keyrepo string // key repository
}

const (
	put = "put"
	get = "get"
)

func main()  {

	env := Env{
		qmgr:     "qm20",
		conname: "229.196.94.34.bc.googleusercontent.com(1414)",
		channel: "EPN.SVRCONN",
		// alter channel(epn.svrconn) chltype(svrconn) sslciph(TLS_RSA_WITH_AES_128_CBC_SHA256)
		cipherspec: "TLS_RSA_WITH_AES_128_CBC_SHA256",
		username: "karson",
		password: "password",
		qname: "Q.A",
		keyrepo: "/home/simon/dev/mq-operator/samples/keystore/key.db",
	}

	qmgr, err := CreateConnection(env)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	qobj, err := OpenQueue(env, qmgr, put)
	if err != nil {
		_ = Disconnect(qmgr)
		log.Fatalf("%v\n", err)
	}

	err = PutMessage(qobj)
	if err != nil {
		_ = CloseQueue(qobj)
		_ = Disconnect(qmgr)
		log.Fatalf("%v\n", err)
	}

	_ = CloseQueue(qobj)
	_ = Disconnect(qmgr)
}

func CreateConnection(env Env) (ibmmq.MQQueueManager, error) {

	// connecting to queue manager
	// queue manager stateful set is exposed as a node-port
	// node port will be available on every node: np1, np2, ...
	// cnt -> np1 -> pod
	// cnt -> lb[np1,np2,..] -> pod

	// with ssl, ssl termination is at the pod
	// pod can be reached with node port, then
	// subj-alt-names: node1.example.com, node2.example.com, node3.example.com
	// subj-alt-names: lb.example.com

	cno := ibmmq.NewMQCNO()

	csp := ibmmq.NewMQCSP()
	csp.AuthenticationType = ibmmq.MQCSP_AUTH_USER_ID_AND_PWD
	csp.UserId = env.username
	csp.Password = env.password

	cno.SecurityParms = csp

	cd := ibmmq.NewMQCD()

	// Fill in required fields in the MQCD channel definition structure
	cd.ChannelName = env.channel
	cd.ConnectionName = env.conname

	log.Printf("Connecting to %s ", cd.ConnectionName)

	if env.keyrepo != "" {
		log.Println("Running in TLS Mode")

		cd.SSLCipherSpec = env.cipherspec

		// The ClientAuth field says whether or not the client needs to present its own certificate
		// Must match the SVRCONN definition.
		// alter channel(epn.svrconn) chltype(svrconn) sslcauth(optional)
		cd.SSLClientAuth = ibmmq.MQSCA_OPTIONAL
	}

	// Reference the CD structure from the CNO
	cno.ClientConn = cd

	if env.keyrepo != "" {
		log.Println("Key Repository has been specified")

		sco := ibmmq.NewMQSCO()
		sco.KeyRepository = env.keyrepo
		cno.SSLConfig = sco
	}

	// Indicate that we definitely want to use the client connection method.
	cno.Options = ibmmq.MQCNO_CLIENT_BINDING

	log.Printf("Attempting connection to %s", env.qmgr)
	qMgr, err := ibmmq.Connx(env.qmgr, cno)

	if err == nil {
		log.Println("Connection succeeded")

	} else {
		log.Printf("%v\n", err)
	}

	return qMgr, err
}

func OpenQueue(env Env, qMgrObject ibmmq.MQQueueManager, action string) (ibmmq.MQObject, error) {

	mqod := ibmmq.NewMQOD()
	openOptions := ibmmq.MQOO_OUTPUT

	switch action {
	case put:
		mqod.ObjectType = ibmmq.MQOT_Q
		mqod.ObjectName = env.qname

	case get:
		openOptions = ibmmq.MQOO_INPUT_SHARED
		mqod.ObjectType = ibmmq.MQOT_Q
		mqod.ObjectName = env.qname
	}

	log.Printf("Attempting open queue/topic %s", env.qname)

	qObject, err := qMgrObject.Open(mqod, openOptions)
	if err != nil {
		log.Printf("%v\n", err)

	} else {
		log.Println("Opened queue", qObject.Name)
	}

	return qObject, err
}

func PutMessage(qobj ibmmq.MQObject) error {

	putmqmd := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()

	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT

	putmqmd.Format = ibmmq.MQFMT_STRING

	msgData := "Hello from Go at " + time.Now().Format(time.RFC3339)

	buffer := []byte(msgData)

	// Now put the message to the queue
	err := qobj.Put(putmqmd, pmo, buffer)

	if err != nil {
		fmt.Println(err)

	} else {
		fmt.Println("Put message to", strings.TrimSpace(qobj.Name))

		// Print the MsgId so it can be used as a parameter to amqsget
		fmt.Println("MsgId:" + hex.EncodeToString(putmqmd.MsgId))
	}

	//mqret := 0
	//if err != nil {
	//	mqret = int((err.(*ibmmq.MQReturn)).MQCC)
	//}

	return err
}

// Disconnect from the queue manager
func Disconnect(qMgrObject ibmmq.MQQueueManager) error {
	err := qMgrObject.Disc()
	if err == nil {
		fmt.Printf("Disconnected from queue manager %s\n", qMgrObject.Name)
	} else {
		fmt.Println(err)
	}
	return err
}

// CloseQueue Close the queue if it was opened
func CloseQueue(object ibmmq.MQObject) error {
	err := object.Close(0)
	if err == nil {
		fmt.Println("Closed queue")
	} else {
		fmt.Println(err)
	}
	return err
}
