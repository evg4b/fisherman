package runner

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fisherman/infrastructure/git"
	"fisherman/utils"
	"os"
)

func (runner *Runner) createContext(configInfo *config.LoadInfo, appPath string) (*context.CommandContext, error) {
	cwd, err := os.Getwd()
	utils.HandleCriticalError(err)

	info, err := git.GetRepositoryInfo(cwd)
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
		Logger:       runner.logger,
	})
	return context, nil
}
