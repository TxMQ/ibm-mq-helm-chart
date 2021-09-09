package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

var logchan chan string

func formatmsg(msg interface{}, skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	fname := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("[%s:%d] %v", filepath.Base(fname), line, msg)
}

func Logmsg(msg interface{}) {
	// format message, skip formatmsg and logmsg
	msg2 := formatmsg(msg, 2)

	if logchan != nil {
		// write to log channel
		logchan <- msg2
	} else {
		// write to std log output
		log.Printf("%v\n", msg2)
	}
}

func Runlogger() {
	// create log channel
	logchan = make(chan string)
	go runlogger()
}

func runlogger() {
	// skip formatmsg
	log.Printf("%s\n", formatmsg("logger started", 1))

	// read log channel
	for msg := range logchan {
		log.Printf("%s\n",  msg)
	}
}