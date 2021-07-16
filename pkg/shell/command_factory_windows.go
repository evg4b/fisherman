package shell

const LineBreak = "\r\n"

const (
	PowerShell = "powershell"
	Cmd        = "cmd"
)

var PlatformDefaultShell = Cmd

var ShellConfigurations = map[string]wrapConfiguration{
	PowerShell: {
		Path:        PowerShell,
		Args:        []string{"-NoProfile", "-NonInteractive", "-NoLogo"},
		PostCommand: "if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }",
	},
	Cmd: {
		Path:        Cmd,
		Args:        []string{"/Q", "/D", "/K"},
		PostCommand: "IF %ERRORLEVEL% NEQ 0 ( exit %ERRORLEVEL% )",
	},
}
