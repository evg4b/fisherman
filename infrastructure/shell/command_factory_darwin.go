package shell

import (
	"os/exec"
	"strings"
)

const LineBreak = "\n"
const PathVariableSeparator = ":"

func init() {
	CommandFactory = func(commands []string) (*exec.Cmd, error) {
		ps, err := exec.LookPath("bash")
		if err != nil {
			return nil, err
		}

		command := strings.Join(commands, LineBreak)

		return exec.Command(ps, "-c", command), nil
	}
}
