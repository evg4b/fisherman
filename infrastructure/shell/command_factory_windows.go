package shell

import (
	"context"
	"os/exec"
	"strings"
)

const LineBreak = "\r\n"
const PathVariableSeparator = ";"

func CommandFactory(ctx context.Context, commands []string) (*exec.Cmd, error) {
	powerShell, err := exec.LookPath("powershell")
	if err != nil {
		return nil, err
	}
	command := strings.Join(commands, LineBreak)

	return exec.CommandContext(ctx, powerShell, "-NoProfile", "-NonInteractive", "-NoLogo", command), nil
}
