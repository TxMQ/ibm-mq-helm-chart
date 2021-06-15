package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

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
	var ctl chan int
	ctl = make(chan int)

	var sig chan os.Signal
	sig = make(chan os.Signal)

	var cld chan os.Signal
	cld = make(chan os.Signal)
	signal.Notify(cld, syscall.SIGCHLD)

	go func() {
		for {
			select {
			case <- cld:
				fmt.Println("zombie...")
			case <- sig:
				fmt.Println("signal, exiting...")
				ctl <- 1
				break
			}
		}
	}()

	select {
	case <- ctl:
		fmt.Println("mqrunner exiting...")
	}
}
