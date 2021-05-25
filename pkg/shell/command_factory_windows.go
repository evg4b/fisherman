package shell

const LineBreak = "\r\n"

const (
	PowerShell = "powershell"
	Cmd        = "cmd"
)

var DefaultShell = Cmd

var ArgumentBuilders = map[string]ArgumentBuilder{
	PowerShell: func() []string { return []string{"-NoProfile", "-NonInteractive", "-NoLogo"} },
	Cmd:        func() []string { return []string{"/Q", "/D"} },
}
