package configuration

import (
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"fisherman/pkg/shell" // TODO: remove this dependency
	"os/user"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const gitDir = ".git"

func FindConfigFiles(usr *user.User, cwd string, fs billy.Filesystem) (map[string]string, error) {
	configs := map[string]string{}
	for _, mode := range ModeOptions {
		folder := GetConfigFolder(usr, cwd, mode)
		files := []string{}
		for _, name := range constants.AppConfigNames {
			configPath := filepath.Join(folder, name)
			exist, err := utils.Exists(fs, configPath)
			if err != nil {
				return nil, err
			}

			if exist {
				files = append(files, configPath)
			}
		}

		if len(files) > 1 {
			return nil, errors.Errorf("more then one config file specifies in folder '%s'", folder)
		}

		if len(files) == 1 {
			configs[mode] = files[0]
		}
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

func Load(fs billy.Filesystem, files map[string]string) (*FishermanConfig, error) {
	config := FishermanConfig{
		Output:       log.DefaultOutputConfig,
		DefaultShell: shell.Default().GetName(),
	}

	for _, mode := range ModeOptions {
		file, ok := files[mode]
		if ok {
			loadedConfig, err := unmarshlFile(fs, file)
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

func unmarshlFile(fs billy.Filesystem, path string) (*FishermanConfig, error) {
	var config FishermanConfig

	file, err := fs.Open(path)
	if err != nil {
		return nil, errors.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
