package initialize

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/hashicorp/go-multierror"
)

type Command struct {
	flagSet  *flag.FlagSet
	mode     string
	force    bool
	absolute bool
	usage    string
	files    billy.Filesystem
	app      internal.AppInfo
	user     *user.User
}

// TODO: Refactor to implement options pattern.
func NewCommand(files billy.Filesystem, app internal.AppInfo, user *user.User) *Command {
	command := &Command{
		flagSet: flag.NewFlagSet("init", flag.ExitOnError),
		usage:   "initializes fisherman in git repository",
		files:   files,
		app:     app,
		user:    user,
	}

	modeMessage := fmt.Sprintf(
		"config location: %s, %s (default), %s. read more here %s",
		configuration.LocalMode,
		configuration.RepoMode,
		configuration.GlobalMode,
		constants.ConfigurationInheritanceURL)

	command.flagSet.StringVar(&command.mode, "mode", configuration.RepoMode, modeMessage)
	command.flagSet.BoolVar(&command.force, "force", false, "forces overwrites existing hooks")
	command.flagSet.BoolVar(&command.absolute, "absolute", false, "used absolute path to binary in hook")

	return command
}

func (c *Command) Run(ctx context.Context, args []string) error {
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	log.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(c.app.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			exist, err := utils.Exists(c.files, hookPath)
			if err != nil {
				return err
			}

			if exist {
				log.Debugf("Hook '%s' already exist", hookName)
				result = multierror.Append(result, errors.Errorf("file %s already exists", hookPath))
			}
		}

		err := result.ErrorOrNil()
		if err != nil {
			return err
		}
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(c.app.Cwd, ".git", "hooks", hookName)

		err := util.WriteFile(c.files, hookPath, buildHook(hookName, c.getBinaryPath(), c.absolute), os.ModePerm)
		if err != nil {
			return err
		}

		log.Infof("Hook '%s' (%s) was writted", hookName, hookPath)
	}

	return c.writeConfig()
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

func (c *Command) writeConfig() error {
	configFolder, err := configuration.GetConfigFolder(c.user, c.app.Cwd, c.mode)
	if err != nil {
		return err
	}

	configPath := filepath.Join(configFolder, constants.AppConfigNames[0])
	exist, err := utils.Exists(c.files, configPath)
	if err != nil {
		return err
	}

	if !exist {
		content := configuration.DefaultConfig
		utils.FillTemplate(&content, map[string]interface{}{
			"URL": constants.ConfigurationFilesDocksURL,
		})

		err := util.WriteFile(c.files, configPath, []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) getBinaryPath() string {
	if c.absolute {
		return c.app.Executable
	}

	return constants.AppName
}
