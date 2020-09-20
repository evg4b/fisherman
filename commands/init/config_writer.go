package init

import (
	"fisherman/config"
	"fisherman/infrastructure/io"

	"gopkg.in/yaml.v2"
)

func writeFishermanConfig(accessor io.FileAccessor, configPath string) error {
	if !accessor.Exist(configPath) {
		content, err := yaml.Marshal(config.DefaultConfig)
		if err != nil {
			return err
		}

		err = accessor.Write(configPath, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
