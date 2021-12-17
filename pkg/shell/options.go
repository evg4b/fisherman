package shell

import (
	"io"

	"golang.org/x/text/encoding"
)

type hostOption = func(str ShellStrategy, host *Host)

// WithStdout setups Stdout writer for shell host.
func WithStdout(output io.Writer) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Stdout = output
	}
}

// WithStderr setups Stderr writer for shell host.
func WithStderr(output io.Writer) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Stderr = output
	}
}

// WithEnv setups environment variables (ShellStrategy defined variables will be included) for shell host.
func WithEnv(env []string) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Env = strategy.EnvWrapper(env)
	}
}

// WithRawEnv setups environment variables (ShellStrategy defined variables will not be included) for shell host.
func WithRawEnv(env []string) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Env = env
	}
}

// WithArgs setups arguments (ShellStrategy defined arguments will be included) for shell host.
func WithArgs(args []string) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Args = strategy.ArgsWrapper(args)
	}
}

// WithRawArgs setups arguments (ShellStrategy defined arguments will not be included) for shell host.
func WithRawArgs(args []string) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Args = args
	}
}

// WithCwd setups environment variables for shell host.
func WithCwd(dir string) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		host.command.Dir = dir
	}
}

// WithEncoding setups shell input/output encoding.
func WithEncoding(encoding encoding.Encoding) hostOption {
	return func(strategy ShellStrategy, host *Host) {
		if encoding != nil {
			host.encoding = encoding
		}
	}
}
