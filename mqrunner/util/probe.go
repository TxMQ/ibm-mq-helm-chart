package util

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Probe struct {
	srv *http.Server
}

func (p *Probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/ready" {
		// todo
		fmt.Println("ready probe called... return 200")

	} else if r.RequestURI == "/healthy" {
		// todo
		fmt.Println("healthy probe called... return 200")

	} else {
		fmt.Println("probe called... return 200")
	}

	w.WriteHeader(200)
	w.Write([]byte("ok"))
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
		fmt.Println("probe running...")
		ctl <- 0

		err := probe.srv.ListenAndServe()
		fmt.Printf("probe stopped... %v\n", err)
	}()

	return probe
}
