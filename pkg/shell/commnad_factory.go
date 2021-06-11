package shell

import (
	"context"
	"os/exec"

	"github.com/go-errors/errors"
)

type ArgumentBuilder = func() []string

func CommandFactory(ctx context.Context, shell string) (*exec.Cmd, error) {
	if builder, ok := ArgumentBuilders[shell]; ok {
		binPath, err := exec.LookPath(shell)
		if err != nil {
			return nil, err
		}

		return exec.CommandContext(ctx, binPath, builder()...), nil
	}

	return nil, errors.Errorf("shell '%s' is not supported", shell)
}
