package util

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//var _stopctl chan int

func StartRunner() chan int {

	// output control channel
	var ctl chan int
	ctl = make(chan int)

	//_stopctl = make(chan int)

	// start probe
	probe := StartProbe(ctl)
	<- ctl

	// setup signals
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
				log.Printf("zobmie...")
				var ws unix.WaitStatus

				pid, err := unix.Wait4(-1, &ws, unix.WNOHANG, nil)
				if err != nil {
					log.Printf("%v\n", err)
				} else {
					log.Printf("Reaped PID %v", pid)
				}

			case <- sig:
				fmt.Println("signal, exiting...")
				// shutdown probe
				_ = probe.Shutdown()

				ctl <- 1
				break

			//case <- _stopctl:
			//	fmt.Println("shutdown, exiting...")
			//	// shutdown probe
			//	_ = probe.Shutdown()
			//
			//	break
			}
		}
	}()

	return ctl
}

//func StopRunner() {
//	_stopctl <- 1
//}