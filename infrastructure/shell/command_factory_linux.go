package shell

import (
	"context"
	"os/exec"
	"strings"
)

const LineBreak = "\n"
const PathVariableSeparator = ":"

func CommandFactory(ctx context.Context, commands []string) (*exec.Cmd, error) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		return nil, err
	}

	command := strings.Join(commands, LineBreak)

	return exec.CommandContext(ctx, bash, "-c", command), nil
}
