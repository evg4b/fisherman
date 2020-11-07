package runner

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/initialize"
	"fisherman/commands/remove"
	"fisherman/commands/version"
	"fisherman/config"
	"fisherman/infrastructure/fs"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
	"fisherman/infrastructure/vcs"
	"flag"
	"os"
	"os/user"
)

func GetRunnerArgs() (Args, error) {
	usr, err := user.Current()
	if err != nil {
		return Args{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return Args{}, err
	}

	appPath, err := os.Executable()
	if err != nil {
		return Args{}, err
	}

	fileAccessor := fs.NewAccessor()

	conf, configInfo, err := config.Load(cwd, usr, fileAccessor)
	if err != nil {
		return Args{}, err
	}

	repo, err := vcs.NewGitRepository(cwd)
	if err != nil {
		return Args{}, err
	}

	log.Configure(conf.Output)

	handling := flag.ExitOnError

	return Args{
		CommandList: []commands.CliCommand{
			initialize.NewCommand(handling),
			handle.NewCommand(handling),
			remove.NewCommand(handling),
			version.NewCommand(handling),
		},
		Config:     conf,
		ConfigInfo: configInfo,
		Files:      fileAccessor,
		SystemUser: usr,
		Cwd:        cwd,
		Executable: appPath,
		Repository: repo,
		Shell:      shell.NewShell(),
	}, nil
}
