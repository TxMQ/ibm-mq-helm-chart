package probe

import (
	"fmt"
	"net/http"
	"szesto.com/mqrunner/logger"
	"szesto.com/mqrunner/util"
)

type probe struct {
	address string
	qmgr string
}

func qmgrReady(qmgr string) bool {

	if running, err := util.IsQmgrRunning(qmgr, true); err != nil {
		return false

	} else if running {
		return true

	} else {
		return false
	}
}

func qmgrHealthy(qmgr string) bool {
	return true
}

func qmgrStarted(qmgr string) bool {
	return true
}

func (p probe) readyProbe(w http.ResponseWriter, req *http.Request) {
	statusCode, resp := serviceOK(p.qmgr)
	if ! qmgrReady(p.qmgr) {
		statusCode, resp = serviceUnavailable(p.qmgr)
	}
	writeResponse(w, statusCode, resp)
}

func (p probe) healthyProbe(w http.ResponseWriter, req *http.Request) {
	statusCode, resp := serviceOK(p.qmgr)
	if ! qmgrHealthy(p.qmgr) {
		statusCode, resp = serviceUnavailable(p.qmgr)
	}
	writeResponse(w, statusCode, resp)
}

func (p probe) startedProbe(w http.ResponseWriter, req *http.Request) {
	statusCode, resp := serviceOK(p.qmgr)
	if ! qmgrStarted(p.qmgr) {
		statusCode, resp = serviceUnavailable(p.qmgr)
	}
	writeResponse(w, statusCode, resp)
}

func serviceOK(qmgr string) (int, string) {
	return http.StatusOK, fmt.Sprintf("qmgr %s ok", qmgr)
}

func serviceUnavailable(qmgr string) (int, string) {
	return http.StatusServiceUnavailable, fmt.Sprintf("qmgr %s unavailable", qmgr)
}

func writeResponse(w http.ResponseWriter, statusCode int, response string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(response))
}

func StartProbes(qmgr string) {
	logger.Logmsg(fmt.Sprintf("starting probes for qmgr '%s'", qmgr))
	go runprobes(qmgr)
}

func runprobes(qmgr string) {
	p := probe{
		address: ":40000",
		qmgr: qmgr,
	}

	// each handler is invoked in a new goroute
	// use default request multiplexer
	http.HandleFunc("/ready", p.readyProbe)
	http.HandleFunc("/healthy", p.healthyProbe)
	http.HandleFunc("/started", p.startedProbe)

	if err := http.ListenAndServe(p.address, nil); err != nil {
		logger.Logmsg(fmt.Sprintf("qmgr probes stopped, %v", err))
	}
}