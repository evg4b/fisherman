package shell

import (
	"context"
	"errors"
	"fisherman/internal/utils"
	"io"
	"os/exec"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// ShellStrategy is interface to describe base concrete shell command.
type ShellStrategy interface { // nolint: revive
	GetName() string
	GetCommand(context.Context) *exec.Cmd
	ArgsWrapper([]string) []string
	EnvWrapper([]string) []string
	GetEncoding() encoding.Encoding
}

// Host is shell host structure to comunicate with shell process.
type Host struct {
	command    *exec.Cmd
	stdin      io.Writer
	closeStdin func() error
	closeOnce  sync.Once
	mu         sync.Mutex
	encoding   encoding.Encoding
}

// NewHost creates new shell host based on passed strategy.
func NewHost(ctx context.Context, strategy ShellStrategy, options ...hostOption) *Host {
	host := &Host{
		command:    strategy.GetCommand(ctx),
		encoding:   strategy.GetEncoding(),
		mu:         sync.Mutex{},
		closeStdin: func() error { return nil },
	}

	for _, option := range options {
		option(strategy, host)
	}

	return host
}

// Write is io.Writer interface implementation.
// Write automatically starts shell process if it has not been started before.
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
	if utils.IsEmpty(script) {
		return errors.New("script can not be empty")
	}

	if err := host.Start(); err != nil {
		return err
	}

	if script[len(script)-1] != '\n' {
		script += "\n"
	}

	if _, err := host.Write([]byte(script)); err != nil {
		return err
	}

	return host.Wait()
}

// Start starts shell host.
func (host *Host) Start() error {
	host.mu.Lock()
	defer host.mu.Unlock()

	if host.isStarted() {
		return errors.New("host already started")
	}

	return host.startUnsave()
}

// Wait waits for the shell to exit and waits for any copying from stdout or stderr to complete.
// Wait automatically closes stdin pipe, and writing will be unavailable after the call.
func (host *Host) Wait() error {
	if err := host.Close(); err != nil {
		return err
	}

	return host.command.Wait()
}

func (host *Host) Close() (err error) {
	host.closeOnce.Do(func() {
		err = host.closeStdin()
	})

	return err
}

func (host *Host) Terminate() error {
	if !host.isStarted() {
		return errors.New("host is not started")
	}

	return host.command.Process.Kill()
}

func (host *Host) isStarted() bool {
	return host.command.Process != nil
}

func (host *Host) startUnsave() error {
	originStdin, err := host.command.StdinPipe()
	if err != nil {
		return err
	}

	decoder := host.encoding.NewDecoder()
	encoder := host.encoding.NewEncoder()

	host.command.Stdout = wrapWriter(host.command.Stdout, decoder)
	host.command.Stderr = wrapWriter(host.command.Stderr, decoder)
	stdin := wrapWriter(originStdin, encoder)

	host.stdin = stdin
	host.closeStdin = func() error {
		encodingWrapperErr := stdin.Close()
		stdinErr := originStdin.Close()

		if encodingWrapperErr != nil {
			return encodingWrapperErr
		}

		return stdinErr
	}

	return host.command.Start()
}

func wrapWriter(w io.Writer, t transform.Transformer) io.WriteCloser {
	if w != nil {
		return transform.NewWriter(w, t)
	}

	return nil
}
