package shell

import (
	"time"

	"golang.org/x/net/context"
)

type ExecResult struct {
	Name     string
	ExitCode int
	Error    error
	Time     time.Duration
}

func (r *ExecResult) IsSuccessful() bool {
	return r.Error == nil && r.ExitCode == 0
}

func (r *ExecResult) IsCanceled() bool {
	return r.Error == context.Canceled
}
