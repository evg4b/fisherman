package runner

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fisherman/infrastructure/git"
	"os"
)

func (runner *Runner) createContext(appPath string) (*context.CommandContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	info, err := git.GetRepositoryInfo(cwd)
	if err != nil {
		return nil, err
	}
	configInfo, err := config.LoadConfig(cwd, runner.systemUser, runner.fileAccessor)
	if err != nil {
		return nil, err
	}
	context := context.NewContext(context.CliCommandContextParams{
		RepoInfo:     info,
		FileAccessor: runner.fileAccessor,
		Usr:          runner.systemUser,
		Cwd:          cwd,
		AppPath:      appPath,
		ConfigInfo:   configInfo,
	})
	return context, nil
}
