package remove

import (
	"fisherman/config"
	"fisherman/constants"
	i "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"flag"
	"os/user"
	"path/filepath"
)

type Command struct {
	flagSet *flag.FlagSet
	usage   string
	files   i.FileSystem
	app     *internal.AppInfo
	user    *user.User
}

func NewCommand(files i.FileSystem, app *internal.AppInfo, user *user.User) *Command {
	return &Command{
		flagSet: flag.NewFlagSet("remove", flag.ExitOnError),
		usage:   "removes fisherman from git repository",
		files:   files,
		app:     app,
		user:    user,
	}
}

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Command) Run() error {
	filesToDelete := []string{
		config.BuildFileConfigPath(c.app.Cwd, c.user, config.RepoMode),
		config.BuildFileConfigPath(c.app.Cwd, c.user, config.LocalMode),
	}

	for _, hookName := range constants.HooksNames {
		filesToDelete = append(filesToDelete, filepath.Join(c.app.Cwd, ".git", "hooks", hookName))
	}

	for _, hookPath := range filesToDelete {
		if c.files.Exist(hookPath) {
			err := c.files.Delete(hookPath)
			if err != nil {
				return err
			}

			log.Infof("File '%s' was removed", hookPath)
		}
	}

	return nil
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}
