package shell

const LineBreak = "\n"

const (
	Bash = "bash"
)

var PlatformDefaultShell = Bash

var ShellConfigurations = map[string]WrapConfiguration{
	Bash: {
		Path: Bash,
		Args: []string{"-i"},
		Init: "set -e",
	},
}
