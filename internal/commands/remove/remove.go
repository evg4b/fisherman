package remove

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"flag"
	"os/user"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

type Command struct {
	flagSet *flag.FlagSet
	usage   string
	files   billy.Filesystem
	app     internal.AppInfo
	user    *user.User
}

// TODO: Refactor to implement options pattern.
func NewCommand(files billy.Filesystem, app internal.AppInfo, user *user.User) *Command {
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

func (command *Command) Run(ctx context.Context) error {
	filesToDelete := []string{}
	for _, config := range command.app.Configs {
		filesToDelete = append(filesToDelete, config)
	}

	for _, hookName := range constants.HooksNames {
		path := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)
		exist, err := utils.Exists(command.files, path)
		if err != nil {
			return err
		}

		if exist {
			filesToDelete = append(filesToDelete, path)
		}
	}

	for _, hookPath := range filesToDelete {
		err := command.files.Remove(hookPath)
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
