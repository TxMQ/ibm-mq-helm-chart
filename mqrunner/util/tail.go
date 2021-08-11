package util

import (
	"fmt"
	"github.com/nxadm/tail"
	"strings"
)

func TailMqLog() {
	// all current error messages are written to this log file
	// periodically this file is copied over to 02.log and truncated
	mqlog := "/var/mqm/errors/AMQERR01.LOG"

	config := tail.Config{
		Location:    nil,
		ReOpen:      true,	// Reopen recreated files (tail -F)
		MustExist:   false,	// Fail early if the file does not exist
		Poll:        false,
		Pipe:        false,
		Follow:      true,	// Continue looking for new lines (tail -f)
		MaxLineSize: 0,
		RateLimiter: nil,
		Logger:      nil,
	}

	go TailFile(mqlog, config)
}

func TailQmgrLog(qmgr string) {
	// all current error messages are written to this log file
	// periodically this file is copied over to 02.log and truncated
	qmgrlog := fmt.Sprintf("/var/mqm/qmgrs/%s/errors/AMQERR01.LOG", qmgr)

	config := tail.Config{
		Location:    nil, 	// Tail from this location. If nil, start at the beginning of the file
		ReOpen:      true, 	// Reopen recreated files (tail -F)
		MustExist:   false, // Fail early if the file does not exist
		Poll:        false, // Poll for file changes instead of using the default inotify
		Pipe:        false, // The file is a named pipe (mkfifo)
		Follow:      true, // Continue looking for new lines (tail -f)
		MaxLineSize: 0, 	// If non-zero, split longer lines into multiple lines
		RateLimiter: nil,	// Use a ratelimiter (e.g. created by the ratelimiter/NewLeakyBucket function)
		Logger:      nil,	// Use a Logger. When nil, the Logger is set to tail.DefaultLogger.
	}

	go TailFile(qmgrlog, config)
}

func TailFile(file string, config tail.Config) {

	// open file
	t, err := tail.TailFile(file, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// range over the channel of Lines
	// forever (follow, re-open)
	for line := range t.Lines {
		// apply AMQ* filter
		if len(line.Text) > 0 {
			if strings.HasPrefix(line.Text, "NOOUT" /*"AMQ"*/) {
				fmt.Println(line.Text)
			}
		}
	}

	// check for error
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
