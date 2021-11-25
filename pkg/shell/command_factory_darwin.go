package shell

const LineBreak = "\n"

const (
	BashAlias = "bash"
)

var PlatformDefaultShell = Bash

var ShellConfigurations = map[string]WrapConfiguration{
	BashAlias: {
		Path: Bash,
		Args: []string{"-i"},
		Init: "set -e",
	},
}
