package shell

import (
	"time"
)

func execWithTime(runFunction func() error) (time.Duration, error) {
	start := time.Now()
	err := runFunction()

	return time.Since(start), err
}
