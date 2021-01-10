package configuration

import (
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"os/user"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const gitDir = ".git"

func Load(cwd string, usr *user.User, files infrastructure.FileSystem) (*FishermanConfig, map[string]string, error) {
	config := FishermanConfig{
		Output: log.DefaultOutputConfig,
	}

	configs := map[string]string{
		GlobalMode: BuildFileConfigPath(cwd, usr, GlobalMode),
		RepoMode:   BuildFileConfigPath(cwd, usr, RepoMode),
		LocalMode:  BuildFileConfigPath(cwd, usr, LocalMode),
	}

	for key, path := range configs {
		if files.Exist(path) {
			log.Debugf("detected %s config file %s", key, path)
			loadedConfig, err := unmarshlFile(path, files)
			if err != nil {
				return nil, nil, err
			}
			err = mergo.Merge(&config, loadedConfig)
			if err != nil {
				return nil, nil, err
			}
		} else {
			configs[key] = ""
		}
	}

	return &config, configs, nil
}

func unmarshlFile(path string, files infrastructure.FileSystem) (*FishermanConfig, error) {
	var config FishermanConfig
	reader, err := files.Reader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func BuildFileConfigPath(cwd string, usr *user.User, mode string) string {
	switch mode {
	case LocalMode:
		return filepath.Join(cwd, gitDir, constants.AppConfigName)
	case RepoMode:
		return filepath.Join(cwd, constants.AppConfigName)
	case GlobalMode:
		return filepath.Join(usr.HomeDir, constants.AppConfigName)
	default:
		panic("unknown config mode")
	}
}
