package remove

import (
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
	filesToDelete := []string{}
	for _, config := range command.app.Configs {
		filesToDelete = append(filesToDelete, config)
	}

	for _, hookName := range constants.HooksNames {
		path := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)
		if command.files.Exist(path) {
			filesToDelete = append(filesToDelete, path)
		}
	}

	for _, hookPath := range filesToDelete {
		err := command.files.Delete(hookPath)
		if err != nil {
			return err
		}

		log.Infof("File '%s' was removed", hookPath)
	}

	return nil
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}
