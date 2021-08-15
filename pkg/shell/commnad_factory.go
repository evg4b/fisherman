package shell

import (
	"os/exec"

	"github.com/go-errors/errors"
)

type WrapConfiguration struct {
	Path        string
	Args        []string
	Init        string
	PostCommand string
	Dispose     string
}

func getShellWrapConfiguration(shell string) (WrapConfiguration, error) {
	if config, ok := ShellConfigurations[shell]; ok {
		binPath, err := exec.LookPath(shell)
		if err != nil {
			return WrapConfiguration{}, err
		}

		config.Path = binPath

		return config, nil
	}

	return WrapConfiguration{}, errors.Errorf("shell '%s' is not supported", shell)
}
