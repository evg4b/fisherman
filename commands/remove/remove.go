package remove

import (
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"flag"
	"os/user"
	"path/filepath"
)

type Command struct {
	flagSet *flag.FlagSet
	usage   string
	files   infrastructure.FileSystem
	app     *internal.AppInfo
	user    *user.User
}

func NewCommand(files infrastructure.FileSystem, app *internal.AppInfo, user *user.User) *Command {
	return &Command{
		flagSet: flag.NewFlagSet("remove", flag.ExitOnError),
		usage:   "removes fisherman from git repository",
		files:   files,
		app:     app,
		user:    user,
	}
}

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run() error {
	filesToDelete := []string{
		configuration.BuildFileConfigPath(command.app.Cwd, command.user, configuration.RepoMode),
		configuration.BuildFileConfigPath(command.app.Cwd, command.user, configuration.LocalMode),
	}

	for _, hookName := range constants.HooksNames {
		filesToDelete = append(filesToDelete, filepath.Join(command.app.Cwd, ".git", "hooks", hookName))
	}

	for _, hookPath := range filesToDelete {
		if command.files.Exist(hookPath) {
			err := command.files.Delete(hookPath)
			if err != nil {
				return err
			}

			log.Infof("File '%s' was removed", hookPath)
		}
	}

	return nil
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}
