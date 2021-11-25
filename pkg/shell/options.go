package shell

import "io"

type (
	shellOption = func(sh *SystemShell)
	hostOption  = func(str ShellStrategy, host *Host)
)

// WithWorkingDirectoryOld setups default working directory for shell.
func WithWorkingDirectoryOld(cwd string) shellOption {
	return func(sh *SystemShell) {
		sh.cwd = cwd
	}
}

// WithDefaultShell setups default shell.
func WithDefaultShell(defaultShell string) shellOption {
	return func(sh *SystemShell) {
		sh.defaultShell = defaultShell
	}
}

// WithEnvOld setups environment variables for shell.
func WithEnvOld(env []string) shellOption {
	return func(sh *SystemShell) {
		sh.env = env
	}
}

// WithStdout setups Stdout writer for shell host.
func WithStdout(output io.Writer) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Stdout = output
	}
}

// WithStderr setups Stderr writer for shell host.
func WithStderr(output io.Writer) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Stderr = output
	}
}

// WithEnv setups environment variables (ShellStrategy defined variables will be included) for shell host.
func WithEnv(env []string) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Env = str.EnvWrapper(env)
	}
}

// WithRawEnv setups environment variables (ShellStrategy defined variables will not be included) for shell host.
func WithRawEnv(env []string) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Env = env
	}
}

// WithArgs setups arguments (ShellStrategy defined arguments will be included) for shell host.
func WithArgs(args []string) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Args = str.ArgsWrapper(args)
	}
}

// WithRawArgs setups arguments (ShellStrategy defined arguments will not be included) for shell host.
func WithRawArgs(args []string) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Args = args
	}
}

// WithCwd setups environment variables for shell host.
func WithCwd(dir string) hostOption {
	return func(str ShellStrategy, host *Host) {
		host.command.Dir = dir
	}
}
