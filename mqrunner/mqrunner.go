package main

import (
	"context"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type probe struct {}
func (p *probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/ready" {
		fmt.Println("ready probe called... return 200")

	} else if r.RequestURI == "/healthy" {
		fmt.Println("healthy probe called... return 200")

	} else {
		fmt.Println("probe called... return 200")
	}

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func main() {

	// setup signals
	var ctl chan int
	ctl = make(chan int)

	srv := &http.Server{
		Addr:           ":40000",
		Handler:        new(probe),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// probe
	go func() {
		fmt.Println("probe running...")
		ctl <- 0

		err := srv.ListenAndServe()
		fmt.Printf("probe stopped... %v\n", err)
	}()

	<- ctl

	var sig chan os.Signal
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	var cld chan os.Signal
	cld = make(chan os.Signal)
	signal.Notify(cld, syscall.SIGCHLD)

	go func() {

		ctl <- 0
		for {
			select {
			case <- cld:
				fmt.Println("zombie...")
				var ws unix.WaitStatus

				pid, err := unix.Wait4(-1, &ws, unix.WNOHANG, nil)
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					fmt.Printf("Reaped PID %v", pid)
				}

			case <- sig:
				fmt.Println("signal, exiting...")
				// shutdown probe
				srv.Shutdown(context.Background())

				ctl <- 1
				break
			}
		}
	}()

	<-ctl // wait for ctl
	fmt.Println("ctl ready...")

	// create mq directories
	cmd := exec.Command("/opt/mqm/bin/crtmqdir", "-f", "-a")
	err := cmd.Run()
	if err != nil {
		// complains about chmod 2775 on /mnt/mqm/data
		fmt.Printf("%v\n", err)
	}

	fmt.Println("directories created...")

	// create queue manager
	cmd = exec.Command("/opt/mqm/bin/crtmqm", "-c", "qm", "-p", "1414", "-q", "qm")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("queue manager created...")

	// start queue manager
	cmd = exec.Command("/opt/mqm/bin/strmqm", "qm")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("queue manager started...")

	// wait for termination
	fmt.Println("waiting for termination...")

	<- ctl
	fmt.Println("mqrunner exiting...")
}
