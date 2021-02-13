package configuration

import (
	"errors"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const gitDir = ".git"

type ConfigLoader struct {
	usr   *user.User
	cwd   string
	files infrastructure.FileSystem
}

func NewLoader(usr *user.User, cwd string, files infrastructure.FileSystem) *ConfigLoader {
	return &ConfigLoader{
		usr:   usr,
		cwd:   cwd,
		files: files,
	}
}

func (loader *ConfigLoader) FindConfigFiles() (map[string]string, error) {
	configs := map[string]string{}

	for _, mode := range []string{GlobalMode, RepoMode, LocalMode} {
		folder := GetConfigFolder(loader.usr, loader.cwd, mode)
		files := []string{}
		for _, name := range constants.AppConfigNames {
			configPath := filepath.Join(folder, name)
			if loader.files.Exist(configPath) {
				files = append(files, configPath)
			}
		}

		if len(files) > 1 {
			return configs, fmt.Errorf("more then one config file specifies in folder '%s'", folder)
		}

		if len(files) == 1 {
			configs[mode] = files[0]
		}
	}

	if len(configs) == 0 {
		return configs, errors.New("no configuration found")
	}

	return configs, nil
}

func GetConfigFolder(usr *user.User, cwd, mode string) string {
	switch mode {
	case LocalMode:
		return filepath.Join(cwd, gitDir)
	case RepoMode:
		return filepath.Join(cwd)
	case GlobalMode:
		return filepath.Join(usr.HomeDir)
	default:
		panic("unknown config mode")
	}
}

func (loader *ConfigLoader) Load(files map[string]string) (*FishermanConfig, error) {
	config := FishermanConfig{
		Output: log.DefaultOutputConfig,
	}

	for _, mode := range []string{GlobalMode, RepoMode, LocalMode} {
		file, ok := files[mode]
		if ok {
			loadedConfig, err := loader.unmarshlFile(file)
			if err != nil {
				return &config, err
			}

			err = mergo.MergeWithOverwrite(&config, loadedConfig)
			if err != nil {
				return &config, err
			}
		}
	}

	return &config, nil
}

func (loader *ConfigLoader) unmarshlFile(path string) (*FishermanConfig, error) {
	var config FishermanConfig
	reader, err := loader.files.Reader(path)
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
