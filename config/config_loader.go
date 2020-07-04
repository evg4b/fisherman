package config

import (
	"fisherman/infrastructure/io"
	"fmt"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFileName = ".fisherman.yaml"

type LoadInfo struct {
	Config           *FishermanConfig
	GlobalConfigPath *string
	RepoConfigPath   *string
	LocalConfigPath  *string
}

func LoadConfig(cwd string, usr *user.User, accessor io.FileAccessor) (*LoadInfo, error) {
	config := FishermanConfig{}

	global, err := unmarshlIfExist(cwd, usr, GlobalMode, accessor, &config)
	if err != nil {
		return nil, err
	}
	repo, err := unmarshlIfExist(cwd, usr, RepoMode, accessor, &config)
	if err != nil {
		return nil, err
	}
	local, err := unmarshlIfExist(cwd, usr, LocalMode, accessor, &config)
	if err != nil {
		return nil, err
	}

	return &LoadInfo{
		Config:           &config,
		GlobalConfigPath: global,
		RepoConfigPath:   repo,
		LocalConfigPath:  local,
	}, nil
}

func unmarshlIfExist(cwd string, usr *user.User, mode string, accessor io.FileAccessor, config *FishermanConfig) (*string, error) {
	path, err := getPathIfExist(cwd, usr, mode, accessor)
	if err != nil {
		return nil, err
	}
	if path != nil {
		data, err := accessor.ReadFile(*path)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal([]byte(data), config)
		if err != nil {
			return nil, err
		}
		return path, nil
	}
	return nil, nil
}

func getPathIfExist(cwd string, usr *user.User, mode string, accessor io.FileAccessor) (*string, error) {
	path, err := BuildFileConfigPath(cwd, usr, mode)
	if err != nil {
		return nil, err
	}
	if accessor.FileExist(path) {
		return &path, nil
	}
	return nil, nil
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
