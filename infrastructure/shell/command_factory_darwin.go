package shell

const LineBreak = "\n"

const (
	Bash = "bash"
)

var DefaultShell = Bash

var ArgumentBuilders = map[string]ArgumentBuilder{
	Bash: func() []string { return []string{"-i"} },
}
