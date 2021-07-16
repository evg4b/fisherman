package shell

import (
	"fmt"
	"os/exec"
)

type wrapConfiguration struct {
	Path        string
	Args        []string
	Init        string
	PostCommand string
	Dispose     string
}

func getShellWrapConfiguration(shell string) (wrapConfiguration, error) {
	if config, ok := ShellConfigurations[shell]; ok {
		binPath, err := exec.LookPath(shell)
		if err != nil {
			return wrapConfiguration{}, err
		}

		config.Path = binPath

		return config, nil
	}

	return wrapConfiguration{}, fmt.Errorf("shell '%s' is not supported", shell)
}
