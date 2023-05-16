package remove

import (
	"context"
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/guards"
	"fisherman/pkg/log"
	"flag"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

type Command struct {
	flagSet     *flag.FlagSet
	usage       string
	fs          billy.Filesystem
	cwd         string
	configFiles map[string]string
}

func NewCommand(options ...removeOption) *Command {
	command := &Command{
		flagSet:     flag.NewFlagSet("remove", flag.ExitOnError),
		usage:       "removes fisherman from git repository",
		configFiles: map[string]string{},
	}

	for _, option := range options {
		option(command)
	}

	guards.ShouldBeDefined(command.fs, "FileSystem should be configured")
	guards.ShouldBeNotEmpty(command.cwd, "Cwd should be configured")

	return command
}

func (c *Command) Run(_ context.Context, args []string) error {
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	filesToDelete := []string{}
	for _, config := range c.configFiles {
		filesToDelete = append(filesToDelete, config)
	}

	for _, hookName := range constants.HooksNames {
		path := filepath.Join(c.cwd, ".git", "hooks", hookName)
		exist, err := utils.Exists(c.fs, path)
		if err != nil {
			return err
		}

		if exist {
			filesToDelete = append(filesToDelete, path)
		}
	}

	for _, hookPath := range filesToDelete {
		err := c.fs.Remove(hookPath)
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
