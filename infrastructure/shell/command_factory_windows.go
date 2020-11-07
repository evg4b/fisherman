package shell

import (
	"os/exec"
	"strings"
)

const LineBreak = "\r\n"
const PathVariableSeparator = ";"

func CommandFactory(commands []string) (*exec.Cmd, error) {
	powerShell, err := exec.LookPath("powershell")
	if err != nil {
		return nil, err
	}
	command := strings.Join(commands, LineBreak)

	return exec.Command(powerShell, "-NoProfile", "-NonInteractive", "-NoLogo", command), nil
}
