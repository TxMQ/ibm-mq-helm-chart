package util

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Probe struct {
	qmgr string
	srv *http.Server
}

var _showready = 5
var _showhealthy = 5

func qmgrReadyProbe(qmgr string, silent bool) bool {

	// if stdby qmgr in multi-instance qmgr, then not ready
	// --- stdby
	// dspmq -x -m qm1
	// QMNAME(qm1) STATUS(Running as standby)
	// INSTANCE(qm1-ibm-mq-0) MODE(Active)
	// INSTANCE(qm1-ibm-mq-1) MODE(Standby)

	if isstandby, err := IsQmgrRunningStandby(qmgr, silent); err != nil {
		// log error
		log.Printf("is-qmgr-ready-probe: %v\n", err)
		return false

	} else if isstandby {
		// running as standby, not ready
		return false
	}

	// --- active
	// QMNAME(qm1) STATUS(Running)
	// INSTANCE(qm1-ibm-mq-0) MODE(Active)
	// INSTANCE(qm1-ibm-mq-1) MODE(Standby)
	//

	if isrunning, err := IsQmgrRunning(qmgr, silent); err != nil {
		// log error
		log.Printf("is-qmgr-ready-probe: %v\n", err)
		return false

	} else if isrunning {
		// running active, ready
		return true
	}

	return false
}

func qmgrHealthyProbe(qmgr string, silent bool) bool {
	return true
}

func (p *Probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	silent := true
	statusCode := http.StatusOK

	if r.RequestURI == "/ready" {

		if _showready > 0 {
			_showready--
			silent = false
		}

		if qmgrReadyProbe(p.qmgr, silent) {
			statusCode = http.StatusOK

			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte("ok"))

		} else {
			statusCode = http.StatusServiceUnavailable

			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte("ServiceUnavailable"))
		}

		if !silent {
			log.Printf("probe: ready probe called, status code %d\n", statusCode)
		}

	} else if r.RequestURI == "/healthy" {

		if _showhealthy > 0 {
			_showhealthy--
			silent = false
		}

		if qmgrHealthyProbe(p.qmgr, silent) {
			statusCode = http.StatusOK

			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte("ok"))

		} else {
			statusCode = http.StatusServiceUnavailable
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte("ServiceUnavailable"))
		}

		if ! silent {
			log.Printf("probe: healthy probe called, status code %d\n", statusCode)
		}

	} else {
		statusCode = http.StatusOK
		log.Printf("probe: unregistered probe called, status code %d\n", statusCode)

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte("ok"))
	}
}

func (p* Probe) Shutdown() error {
	err := p.srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func StartProbe(qmgr string, ctl chan int) *Probe {

	probe := &Probe{}

	// set qmgr on the probe
	probe.qmgr = qmgr

	// configure server
	probe.srv = &http.Server{
		Addr:           ":40000",
		Handler:        probe,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// run probe
	go func() {
		log.Printf("start-probe: %s\n", "probe running...")
		ctl <- 0

		err := probe.srv.ListenAndServe()
		log.Printf("probe stopped... %v\n", err)
	}()

	return probe
}
