package mqsc

import (
	"fmt"
	"strings"
)

type Qalter struct {

	// alter, no output
	//boqname('')
	//bothresh(0)

	// alter, no output
	//defbind(open)
	//defreada(no)
	//defsopt(shared)

	// alter, no output
	//distl(no)
	//nohardenbo
	//imgrcovq(qmgr)

	// alter, no output
	//npmclass(normal)
	//process('')
	//propctl(compat)

	// alter, no output
	//retintvl(999_999_999)
	//defpresp(%s) // default put response sync/async
}

type Qevents struct {
	Qdepthhi int //qdepthhi(80)
	Qdepthlo int //qdepthlo(20)
	Qdepthhiev bool //qdphiev(disabled)
	Qdepthloev bool //qdploev(disabled)
	Qdpmaxev bool //qdpmaxev(enabled)
	Qsvcev bool //qsvciev(none)
	Qsvcint int //qsvcint(999_999_999)
}

type Qtrigger struct {
	Enabled bool //notrigger
	Trigdata string //trigdata('')
	Trigdepth int //trigdpth(1)
	Trigmpri int //trigmpri(0)
	Trigtype string // trigtype(first)
	Initq string // initq('')
}

type Qcluster struct {
	Clusnl string //clusnl('')
	Cluster string //cluster('')
	Clwlprty int //clwlprty(0)
	Clwlrank int //clwlrank(0)
	Clwluseq string //clwluseq(qmgr)
}

type Localqueue struct {

	Name string // queue-name - queue name - param
	Descr string // descr(%s) - plain text comment - param
	Like string // like - name of a queue to model this def - param

	Put string // enable put 	get(%s) enabled/disabled
	Get string // enable get 	put(%s) enabled/disabled

	DefaultPriority int //	defprty(%d)
	DefaultPersistence bool // 	defpsist(%s) yes/no

	Maxdepth int // maxdepth(5000)
	Maxfsize int // maxfsize(2_088_960)
	Maxmsgl int // maxmsgl(4_193_304)
	MsgDeliverySeq string // msgdlvsq(priority)

	Qtrigger Qtrigger
	Qevents Qevents
	Qcluster Qcluster

	acctq string // 	acctq(qmgr)
	monq string // 		monq(qmgr)
	statq string // 	statq(qmgr)
	usage string // 	usage(normal)
	share bool // true 	share

	Authority []Authrec
	Alter []string
}

func (lq *Localqueue) Mqsc() string {

	/*
	DEFINE QLOCAL( q_name )
	   [ ACCTQ( QMGR | ON | OFF ) ]            [ BOQNAME( string ) ]
	   [ BOTHRESH( integer ) ]                 [ CLCHNAME( channel_name ) ]
	   [ CLUSNL( namelist_name ) ]             [ CLUSTER( cluster_name ) ]
	   [ CLWLPRTY( integer ) ]                 [ CLWLRANK( integer ) ]
	   [ CLWLUSEQ( LOCAL | ANY | QMGR ) ]      [ CUSTOM( string ) ]
	   [ DEFBIND( NOTFIXED | OPEN | GROUP ) ]  [ DEFPRESP( SYNC | ASYNC ) ]
	   [ DEFPRTY( integer ) ]                  [ DEFPSIST( YES | NO ) ]
	   [ DEFREADA( NO | YES | DISABLED ) ]     [ DEFSOPT( EXCL | SHARED ) ]
	   [ DESCR( string ) ]                     [ DISTL( YES | NO ) ]
	   [ GET( ENABLED | DISABLED ) ]           [ HARDENBO | NOHARDENBO ]
	   [ IMGRCOVQ( YES | NO | QMGR ) ]         [ INITQ( string ) ]
	   [ LIKE( qlocal_name ) ]                 [ MAXDEPTH( integer ) ]
	   [ MAXFSIZE( DEFAULT | integer ) ]
	   [ MAXMSGL( integer ) ]
	   [ MONQ( OFF | QMGR | LOW | MEDIUM | HIGH ) ]
	   [ MSGDLVSQ( PRIORITY | FIFO ) ]         [ NPMCLASS( NORMAL | HIGH ) ]
	   [ PROCESS( string ) ]
	   [ PROPCTL( COMPAT | NONE | ALL | FORCE | V6COMPAT ) ]
	   [ PUT( ENABLED | DISABLED ) ]           [ QDEPTHHI( integer ) ]
	   [ QDEPTHLO( integer ) ]                 [ QDPHIEV( ENABLED | DISABLED ) ]
	   [ QDPLOEV( ENABLED | DISABLED ) ]       [ QDPMAXEV( ENABLED | DISABLED ) ]
	   [ QSVCIEV( NONE | HIGH | OK ) ]         [ QSVCINT( integer ) ]
	   [ REPLACE | NOREPLACE ]                 [ RETINTVL( integer ) ]
	   [ SCOPE( QMGR | CELL ) ]                [ SHARE | NOSHARE ]
	   [ STATQ( QMGR | ON | OFF ) ]            [ TRIGDATA( string ) ]
	   [ TRIGDPTH( integer ) ]                 [ TRIGGER | NOTRIGGER ]
	   [ TRIGMPRI( integer ) ]
	   [ TRIGTYPE( FIRST | EVERY | DEPTH | NONE ) ]
	   [ USAGE( NORMAL | XMITQ ) ]
	 */

	var mqsc []string

	//if len(lq.Like) > 0 {
	//
	//	t :=
	//		"define qlocal(%s) replace" + cont + // qname
	//		"descr(%s)" + cont + // descr
	//		"like(%s)" // like
	//
	//	descr := fmt.Sprintf("local queue %s", lq.Name)
	//	if len(lq.Descr) > 0 { descr = lq.Descr}
	//
	//	s := fmt.Sprintf(t, lq.Name, descr, lq.Like)
	//
	//	mqsc = append(mqsc, s)
	//

	//t :=
	//	"define qlocal(%s) replace" + cont + // qname
	//	"descr('%s')" + cont + // descr
	//	"put(%s)" + cont + // put enabled/disabled
	//	"get(%s)" + cont + // get enabled/disabled
	//	"defprty(%d)" + cont + // defpriority
	//	"defpsist(%s)" + cont + // default-persistence yes/no
	//	"maxdepth(%d)" + cont + // maxdepth
	//	"maxfsize(%d)" + cont + // maxfsize
	//	"maxmsgl(%d)" + cont + // maxmsgl
	//	"msgdlvsq(%s)" + cont + // msg-delivery-seq (PRIORITY|FIFO)
	//
	//	// monitoring unexported fields
	//	"acctq(%s)" + cont + // acctq(qmgr)
	//	"monq(%s)" + cont + // monq(qmgr)
	//	"statq(%s)" + cont + // statq(qmgr)
	//
	//	"usage(%s)" + cont + // usage(normal)
	//	"%s" + cont + // trigger/notrigger
	//	"share" + endl
	//
	//descr := fmt.Sprintf("local queue %s", lq.Name)
	//if len(lq.Descr) > 0 { descr = lq.Descr}
	//
	//put := "enabled"
	//if len(lq.Put) > 0 && strings.ToLower(lq.Put) == "disabled" {put = "disabled"}
	//
	//get := "enabled"
	//if len(lq.Get) > 0 && strings.ToLower(lq.Get) == "disabled" {get = "disabled"}
	//
	//defpersist := "yes"
	//if lq.DefaultPersistence == false { defpersist = "no" }
	//
	//msgdeliveryseq := "PRIORITY"
	//if strings.ToUpper(lq.MsgDeliverySeq) == "FIFO" {
	//	msgdeliveryseq = "FIFO"
	//}
	//
	//// monitoring
	//mon := "qmgr"
	//
	//// local queue
	//lqnormal := "normal"
	//
	//// trigger
	//trigger := "notrigger"
	//if lq.Qtrigger.Enabled {
	//	// trigger params
	//}
	//
	//s := fmt.Sprintf(t, lq.Name, descr, put, get, lq.DefaultPriority, defpersist,
	//	lq.Maxdepth, lq.Maxfsize, lq.Maxmsgl, msgdeliveryseq,
	//	mon, mon, mon, lqnormal, trigger)
	//
	//mqsc = append(mqsc, s)

	t := "define qlocal(%s) replace" + endl // qname
	s := fmt.Sprintf(t, lq.Name)
	mqsc = append(mqsc, s)

	// authorities
	for _, authrec := range lq.Authority {
		s := authrec.Mqsc(lq.Name, "QUEUE")
		mqsc = append(mqsc, s)
	}

	// alter

	return strings.Join(mqsc, "\n")
}
