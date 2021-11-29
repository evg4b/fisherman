package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// ShellStrategy is interface to describe base concrete shell command.
type ShellStrategy interface { // nolint: revive
	GetCommand(context.Context) *exec.Cmd
	ArgsWrapper([]string) []string
	EnvWrapper([]string) []string
}

// Host is shell host structure to comunicate with shell process.
type Host struct {
	command     *exec.Cmd
	stdin       io.WriteCloser
	stdinClosed bool
	mu          sync.Mutex
}

// NewHost creates new shell host based on passed strategy.
func NewHost(ctx context.Context, shellStr ShellStrategy, options ...hostOption) *Host {
	command := shellStr.GetCommand(ctx)
	host := &Host{
		command: command,
		mu:      sync.Mutex{},
	}

	for _, option := range options {
		option(shellStr, host)
	}

	return host
}

// Write is io.Writer interface implementation.
//
// Write writes len(payload) bytes from payload to shell porecess's Stdin.
// It returns the number of bytes written from payload (0 <= n <= len(payload))
// and any error encountered that caused the write to stop early.
// Write automatically starts shell process if it has not been started before.
// Write returns a non-nil error if it returns n < len(payload).
// Write does not modify the slice data, even temporarily.
// Write does not retain payload.
func (host *Host) Write(payload []byte) (int, error) {
	host.mu.Lock()
	defer host.mu.Unlock()

	if !host.isStarted() {
		if err := host.startUnsave(); err != nil {
			return 0, err
		}
	}

	return host.stdin.Write(payload)
}

// Run runs new shell host based on passed strategy.
func (host *Host) Run(script string) error {
	if err := host.Start(); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(host, script); err != nil {
		return err
	}

	return host.Wait()
}

func (host *Host) Start() error {
	host.mu.Lock()
	defer host.mu.Unlock()

	if !host.isStarted() {
		return host.startUnsave()
	}

	return errors.New("shell host: already started")
}

// Wait waits for the shell to exit and waits for any copying to
// stdin or copying from stdout or stderr to complete.
//
// The command must have been started by Start.
//
// The returned error is nil if the shell runs, has no problems
// copying stdin, stdout, and stderr, and exits with a zero exit
// status.
//
// If the shell script fails to run or doesn't complete successfully, the
// error is of type *ExitError. Other error types may be
// returned for I/O problems.
//
// Wait automatically closes stdin pipe, and writing will be unavailable after the call.
func (host *Host) Wait() error {
	if err := host.Close(); err != nil {
		return err
	}

	return host.command.Wait()
}

func (host *Host) Close() error {
	host.stdinClosed = true

	if host.stdin != nil {
		return host.stdin.Close()
	}

	return nil
}

func (host *Host) isStarted() bool {
	return host.command.Process != nil
}

func (host *Host) startUnsave() error {
	stdin, err := host.command.StdinPipe()
	if err != nil {
		return err
	}

	host.stdin = stdin
	host.stdinClosed = false

	return host.command.Start()
}
