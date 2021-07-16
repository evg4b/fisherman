package shell

const LineBreak = "\n"

const (
	Bash = "bash"
)

var DefaultShell = Bash

var ShellConfigurations = map[string]wrapConfiguration{
	Bash: {
		Path: Bash,
		Args: []string{"-i"},
		Init: "set -e",
	},
}
