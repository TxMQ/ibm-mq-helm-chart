package util

import (
	"fmt"
	"github.com/nxadm/tail"
	"os"
	"strings"
	"szesto.com/mqrunner/logger"
)

func mqErrlogPath() string {
	return "/var/mqm/errors/AMQERR01.LOG"
}

func qmrgErrlogPath(qmgr string) string {
	return fmt.Sprintf("/var/md/%s/errors/AMQERR01.LOG", qmgr)
}

func TailMqLog() {
	// all current error messages are written to this log file
	// periodically this file is copied over to 02.log and truncated
	mqlog := mqErrlogPath()

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

	filters := getLogFilters()
	go TailFile(mqlog, config, filters)
}

func TailQmgrLog(qmgr string) {
	// all current error messages are written to this log file
	// periodically this file is copied over to 02.log and truncated
	qmgrlog := qmrgErrlogPath(qmgr)

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

	filters := getLogFilters()
	go TailFile(qmgrlog, config, filters)
}

func getLogFilters() []string {
	// comma separated list of filters
	fenv := os.Getenv("MQ_LOG_FILTER")
	if len(fenv) > 0 {
		return strings.Split(fenv, ",")
	}
	return nil
}

func TailFile(file string, config tail.Config, filters []string) {

	const mqLogDefaultFilter = "DEFAULT_FILTER"
	const mqLogNoFilter = "NO_FILTER"

	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("applying log filters: %v to log file: %s", filters, file))
	}

	// open file
	t, err := tail.TailFile(file, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// range over the channel of Lines
	// forever (follow, re-open)
	for line := range t.Lines {

		if len(filters) == 0 {
			// no-output
			continue

		} else if len(filters) == 1 && strings.ToUpper(filters[0]) == mqLogNoFilter {
			// output every line
			fmt.Println(line.Text)

		} else if len(filters) == 1 && strings.ToUpper(filters[0]) == mqLogDefaultFilter {
			if strings.HasPrefix(line.Text, "AMQ") {
				fmt.Println(line.Text)
			}

		} else if len(filters) > 0 {
			// apply input filters
			for _, f := range filters {
				if strings.HasPrefix(line.Text, f) {
					fmt.Println(line.Text)
					break
				}
			}
		}
	}

	// check for error
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
