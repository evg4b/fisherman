package shell

import (
	"time"
)

func ExecWithTime(runFunction func() error) (time.Duration, error) {
	start := time.Now()
	err := runFunction()

	return time.Since(start), err
}
