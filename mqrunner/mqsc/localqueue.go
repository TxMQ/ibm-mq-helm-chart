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

	Put bool // enable put 	get(%s) enabled/disabled
	Get bool // enable get 	put(%s) enabled/disabled

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

	var mqsc []string

	if len(lq.Like) > 0 {

		t :=
			"define qlocal(%s) replace" + cont + // qname
			"descr(%s)" + cont + // descr
			"like(%s)" // like

		descr := fmt.Sprintf("local queue %s", lq.Name)
		if len(lq.Descr) > 0 { descr = lq.Descr}

		s := fmt.Sprintf(t, lq.Name, descr, lq.Like)

		mqsc = append(mqsc, s)

	} else {

		t :=
			"define qlocal(%s) replace" + cont + // qname
			"descr('%s')" + cont + // descr
			"put(%s)" + cont + // put enabled/disabled
			"get(%s)" + cont + // get enabled/disabled
			"defprty(%d)" + cont + // defpriority
			"defpsist(%s)" + cont + // default-persistence yes/no
			"maxdepth(%d)" + cont + // maxdepth
			"maxfsize(%d)" + cont + // maxfsize
			"maxmsgl(%d)" + cont + // maxmsgl
			"msgdlvsq(%s)" + cont + // msg-delivery-seq

			// monitoring unexported fields
			"acctq(%s)" + cont + // acctq(qmgr)
			"monq(%s)" + cont + // monq(qmgr)
			"statq(%s)" + cont + // statq(qmgr)

			"usage(%s)" + cont + // usage(normal)
			"%s" + cont + // trigger/notrigger
			"share" + endl

		descr := fmt.Sprintf("local queue %s", lq.Name)
		if len(lq.Descr) > 0 { descr = lq.Descr}

		put := "disabled"
		if lq.Put {put = "enabled"}

		get := "disabled"
		if lq.Get {get = "enabled"}

		defpersist := "yes"
		if lq.DefaultPersistence == false { defpersist = "no" }

		// monitoring
		mon := "qmgr"

		// local queue
		lqnormal := "normal"

		// trigger
		trigger := "notrigger"
		if lq.Qtrigger.Enabled {
			// trigger params
		}

		s := fmt.Sprintf(t, lq.Name, descr, put, get, lq.DefaultPriority, defpersist,
			lq.Maxdepth, lq.Maxfsize, lq.Maxmsgl, lq.MsgDeliverySeq,
			mon, mon, mon, lqnormal, trigger)

		mqsc = append(mqsc, s)
	}

	// alter

	return strings.Join(mqsc, "\n")
}

/*
 * parameter list
 *
boqname - backout reque name - alt
bothresh - backout threshold - alt
capexpry
csusnl - list of clusters for the queue - cluster
cluster - name of the cluster for the queue - cluster
clwlprty - cluster workload distribution pri - cluster
clwlrank - rank of the queue for workload distr - cluster
clwluseq - cluster put behaviour - cluster
custom - custom attribute for new features - alt
defbind - binding for mqbind for the cluster queue - alt
defpresp - put response - alt
defreada - default read-ahead for non-persistent msgs - alt
defsopt - default share option - alt
hardenbo/nohardenbo - alt
imgrcovq - recovery from media image - alt
maxfsize - max size in gb that a queue file can grow (default) - alt
monq - collect monitoring data for queue - qmgr - alt
noreplace
npmclass - reliability for non persistent messages normal/high - alter
process - application to start if trigger is firing - alter
propctl - message property read control - alt
qdepthhi - threshhold for queue-depth-high event - alt
qdepthlo - threshhold for queue-depth-low event - alt
qdphiev - generate hi depth events - alt
qdploev - generate lo depth events - alt
qsvciev - service interval high event - alt
qsvcint - service interval for serivce high events - alt
retintvl - queue no longer needed - alt
scope - queue def scope - alt
 */
