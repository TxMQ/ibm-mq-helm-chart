package util

import (
	"github.com/nxadm/tail"
	"testing"
)

func TestTailFile(t *testing.T) {
	file := "/var/mqm/errors/AMQERR01.LOG"

	config := tail.Config{
		Location:    nil,
		ReOpen:      false, // agree with follow
		MustExist:   false,
		Poll:        false,
		Pipe:        false,
		Follow:      false, // so we can exit
		MaxLineSize: 0,
		RateLimiter: nil,
		Logger:      nil,
	}

	TailFile(file, config)
}
