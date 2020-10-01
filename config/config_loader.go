package config

import (
	"fisherman/constants"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"os/user"
	"path/filepath"

	"errors"

	"gopkg.in/yaml.v3"
)

const gitDir = ".git"

// LoadInfo is
type LoadInfo struct {
	GlobalConfigPath string
	RepoConfigPath   string
	LocalConfigPath  string
}

// LoadConfig is demo
func LoadConfig(cwd string, usr *user.User, accessor inf.FileAccessor) (*FishermanConfig, *LoadInfo) {
	config := FishermanConfig{
		Output: logger.DefaultOutputConfig,
	}

	// Please do not change the order
	loadInfo := &LoadInfo{
		GlobalConfigPath: unmarshlIfExist(cwd, usr, GlobalMode, accessor, &config),
		RepoConfigPath:   unmarshlIfExist(cwd, usr, RepoMode, accessor, &config),
		LocalConfigPath:  unmarshlIfExist(cwd, usr, LocalMode, accessor, &config),
	}

	return &config, loadInfo
}

func unmarshlIfExist(cwd string, usr *user.User, mode string, files inf.FileAccessor, config *FishermanConfig) string {
	path, err := BuildFileConfigPath(cwd, usr, mode)
	utils.HandleCriticalError(err)

	if files.Exist(path) {
		data, err := files.Read(path)
		utils.HandleCriticalError(err)
		err = yaml.Unmarshal([]byte(data), config)
		utils.HandleCriticalError(err)

		return path
	}

	return ""
}

// BuildFileConfigPath returns path to config by config mode
func BuildFileConfigPath(cwd string, usr *user.User, mode string) (string, error) {
	switch mode {
	case LocalMode:
		return filepath.Join(cwd, gitDir, constants.AppConfigName), nil
	case RepoMode:
		return filepath.Join(cwd, constants.AppConfigName), nil
	case GlobalMode:
		return filepath.Join(usr.HomeDir, constants.AppConfigName), nil
	default:
		return "", errors.New("unknown config mode")
	}
}
