package shell

import (
	"os/exec"
	"strings"
)

const windowsLineBreak = "\n\r"

func init() {
	CommandFactory = func(commands []string) (*exec.Cmd, error) {
		ps, err := exec.LookPath("powershell")
		if err != nil {
			return nil, err
		}

		command := strings.Join(commands, windowsLineBreak)

		return exec.Command(ps, "-NoProfile", "-NonInteractive", command), nil
	}
}
