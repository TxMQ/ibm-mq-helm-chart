package util

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Probe struct {
	srv *http.Server
}

var _showready = 5
var _showhealthy = 5

func (p *Probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/ready" {
		// todo
		if  _showready > 0 {
			_showready--
			log.Printf("probe: %s\n", "ready probe called... return 200")
		}

	} else if r.RequestURI == "/healthy" {
		// todo
		if _showhealthy > 0 {
			_showhealthy--
			log.Printf("probe: %s\n", "healthy probe called... return 200")
		}

	} else {
		log.Printf("probe: %s\n", "probe called... return 200")
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}

func (p* Probe) Shutdown() error {
	err := p.srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func StartProbe(ctl chan int) *Probe {

	probe := &Probe{}

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
