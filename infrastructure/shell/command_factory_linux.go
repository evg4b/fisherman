package shell

import (
	"os/exec"
	"strings"
)

const linuxLineBreak = "\n"

func init() {
	CommandFactory = func(commands []string) (*exec.Cmd, error) {
		ps, err := exec.LookPath("bash")
		if err != nil {
			return nil, err
		}

		command := strings.Join(commands, linuxLineBreak)

		return exec.Command(ps, "-c", command), nil
	}
}
