package config

import (
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"fmt"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFileName = ".fisherman.yaml"

// LoadInfo is
type LoadInfo struct {
	GlobalConfigPath string
	RepoConfigPath   string
	LocalConfigPath  string
}

// LoadConfig is demo
func LoadConfig(cwd string, usr *user.User, accessor io.FileAccessor) (*FishermanConfig, *LoadInfo, error) {
	config := FishermanConfig{
		Output: logger.DefaultOutputConfig,
	}

	global, err := unmarshlIfExist(cwd, usr, GlobalMode, accessor, &config)
	if err != nil {
		return nil, nil, err
	}

	repo, err := unmarshlIfExist(cwd, usr, RepoMode, accessor, &config)
	if err != nil {
		return nil, nil, err
	}

	local, err := unmarshlIfExist(cwd, usr, LocalMode, accessor, &config)
	if err != nil {
		return nil, nil, err
	}

	loadInfo := &LoadInfo{
		GlobalConfigPath: global,
		RepoConfigPath:   repo,
		LocalConfigPath:  local,
	}

	return &config, loadInfo, nil
}

func unmarshlIfExist(cwd string, usr *user.User, mode string, accessor io.FileAccessor, config *FishermanConfig) (string, error) {
	path, err := getPathIfExist(cwd, usr, mode, accessor)
	if err != nil {
		return "", err
	}

	if !utils.IsEmpty(path) {
		data, err := accessor.ReadFile(path)
		if err != nil {
			return "", err
		}

		err = yaml.Unmarshal([]byte(data), config)
		if err != nil {
			return "", err
		}

		return path, nil
	}

	return "", nil
}

func getPathIfExist(cwd string, usr *user.User, mode string, accessor io.FileAccessor) (string, error) {
	path, err := BuildFileConfigPath(cwd, usr, mode)
	if err != nil {
		return "", err
	}
	if accessor.FileExist(path) {
		return path, nil
	}
	return "", nil
}

func BuildFileConfigPath(cwd string, usr *user.User, mode string) (string, error) {
	switch mode {
	case "local":
		return filepath.Join(cwd, ".git", configFileName), nil
	case "repo":
		return filepath.Join(cwd, configFileName), nil
	case "global":
		return filepath.Join(usr.HomeDir, configFileName), nil
	default:
		return "", fmt.Errorf("unknow mode")
	}
}
