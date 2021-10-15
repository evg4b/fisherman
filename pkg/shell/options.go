package shell

type shellOption = func(sh *SystemShell)

// WithWorkingDirectory setups working directory for shell.
func WithWorkingDirectory(cwd string) shellOption {
	return func(sh *SystemShell) {
		sh.cwd = cwd
	}
}

// WithWorkingDirectory setups default shell.
func WithDefaultShell(defaultShell string) shellOption {
	return func(sh *SystemShell) {
		sh.defaultShell = defaultShell
	}
}
