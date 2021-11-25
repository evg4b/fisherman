package shell

const LineBreak = "\r\n"

const (
	PowerShellAlias = "powershell"
	CmdAlias        = "cmd"
)

var PlatformDefaultShell = CmdAlias

var ShellConfigurations = map[string]WrapConfiguration{
	PowerShellAlias: {
		Path:        PowerShellAlias,
		Args:        []string{"-NoProfile", "-NonInteractive", "-NoLogo"},
		PostCommand: "if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }",
	},
	CmdAlias: {
		Path:        CmdAlias,
		Args:        []string{"/Q", "/D", "/K"},
		PostCommand: "IF %ERRORLEVEL% NEQ 0 ( exit %ERRORLEVEL% )",
	},
}
