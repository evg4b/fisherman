package init

import (
	"fisherman/config"
	"fisherman/infrastructure/io"
	"gopkg.in/yaml.v2"
	"os/user"
)

func WriteFishermanConfig(cwd string, usr *user.User, mode string, accessor io.FileAccessor) error {
	configPath, err := config.BuildFileConfigPath(cwd, usr, mode)
	if err != nil {
		return err
	}
	if !accessor.FileExist(configPath) {
		content, err := yaml.Marshal(config.DefaultConfig)
		if err != nil {
			return err
		}
		err = accessor.WriteFile(configPath, string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
