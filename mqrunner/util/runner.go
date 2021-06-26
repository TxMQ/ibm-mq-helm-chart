package util

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"syscall"
)

func StartRunner() chan int {

	// control channel
	var ctl chan int
	ctl = make(chan int)

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
				_ = probe.Shutdown()

				ctl <- 1
				break
			}
		}
	}()

	return ctl
}
