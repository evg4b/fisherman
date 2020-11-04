package shell

import (
	"os/exec"
	"strings"
)

const LineBreak = "\r\n"
const PathVariableSeparator = ";"

func CommandFactory(commands []string) (*exec.Cmd, error) {
	ps, err := exec.LookPath("powershell")
	if err != nil {
		return nil, err
	}

	command := strings.Join(commands, LineBreak)

	return exec.Command(ps, "-NoProfile", "-NonInteractive", command), nil
}
