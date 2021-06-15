package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	// setup signals
	var ctl chan int
	ctl = make(chan int)

	var sig chan os.Signal
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

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

	// create queue manager
	cmd = exec.Command("/opt/mqm/bin/crtmqm", "-c", "qm", "-p", "1414", "-q", "qm")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// start queue manager
	cmd = exec.Command("/opt/mqm/bin/strmqm", "qm")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// wait for termination
	<- ctl
	fmt.Println("mqrunner exiting...")
}
