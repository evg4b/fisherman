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

func (c *Command) Run(ctx context.Context, args []string) error {
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	filesToDelete := []string{}
	for _, config := range c.app.Configs {
		filesToDelete = append(filesToDelete, config)
	}

	for _, hookName := range constants.HooksNames {
		path := filepath.Join(c.app.Cwd, ".git", "hooks", hookName)
		exist, err := utils.Exists(c.files, path)
		if err != nil {
			return err
		}

		if exist {
			filesToDelete = append(filesToDelete, path)
		}
	}

	for _, hookPath := range filesToDelete {
		err := c.files.Remove(hookPath)
		if err != nil {
			return err
		}

		log.Infof("File '%s' was removed", hookPath)
	}

	return nil
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}
