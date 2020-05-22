package shell

import (
	"io"

	"golang.org/x/text/encoding"
)

type HostOption = func(str Strategy, host *Host)

// WithStdout setups Stdout writer for shell host.
func WithStdout(output io.Writer) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Stdout = output
	}
}

// WithStderr setups Stderr writer for shell host.
func WithStderr(output io.Writer) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Stderr = output
	}
}

// WithEnv setups environment variables (ShellStrategy defined variables will be included) for shell host.
func WithEnv(env []string) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Env = strategy.EnvWrapper(env)
	}
}

// WithRawEnv setups environment variables (ShellStrategy defined variables will not be included) for shell host.
func WithRawEnv(env []string) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Env = env
	}
}

// WithArgs setups arguments (ShellStrategy defined arguments will be included) for shell host.
func WithArgs(args []string) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Args = strategy.ArgsWrapper(args)
	}
}

// WithRawArgs setups arguments (ShellStrategy defined arguments will not be included) for shell host.
func WithRawArgs(args []string) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Args = args
	}
}

// WithCwd setups environment variables for shell host.
func WithCwd(dir string) HostOption {
	return func(strategy Strategy, host *Host) {
		host.command.Dir = dir
	}
}

// WithEncoding setups shell input/output encoding.
func WithEncoding(encoding encoding.Encoding) HostOption {
	return func(strategy Strategy, host *Host) {
		if encoding != nil {
			host.encoding = encoding
		}
	}
}
