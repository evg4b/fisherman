package shell

import (
	"time"

	"golang.org/x/net/context"
)

type ExecResult struct {
	Name  string
	Error error
	Time  time.Duration
}

func (r *ExecResult) IsSuccessful() bool {
	return r.Error == nil
}

func (r *ExecResult) IsCanceled() bool {
	return r.Error == context.Canceled
}
