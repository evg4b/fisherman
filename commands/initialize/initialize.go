package initialize

import (
	"fisherman/config"
	"fisherman/constants"
	i "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/utils"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type Command struct {
	flagSet *flag.FlagSet
	mode    string
	force   bool
	abslute bool
	usage   string
	files   i.FileSystem
	app     *internal.AppInfo
	user    *user.User
}

func NewCommand(files i.FileSystem, app *internal.AppInfo, user *user.User) *Command {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	command := &Command{
		flagSet: flagSet,
		usage:   "initializes fisherman in git repository",
		files:   files,
		app:     app,
		user:    user,
	}

	modeMessage := fmt.Sprintf(
		"config location (%s, %s (default), %s)",
		config.LocalMode,
		config.RepoMode,
		config.GlobalMode)

	flagSet.StringVar(&command.mode, "mode", config.RepoMode, modeMessage)
	flagSet.BoolVar(&command.force, "force", false, "forces overwrites existing hooks")
	flagSet.BoolVar(&command.abslute, "abs", false, "used absolute path to binary in hook")

	return command
}

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Command) Run() error {
	log.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(c.app.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if c.files.Exist(hookPath) {
				log.Debugf("Hook '%s' already exist", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		return result.ErrorOrNil()
	}

	bin := c.app.Executable
	if !c.abslute {
		bin = utils.NormalizePath(bin)
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(c.app.Cwd, ".git", "hooks", hookName)
		err := c.files.Write(hookPath, buildHook(bin, hookName))
		if err != nil {
			return err
		}
		log.Infof("Hook '%s' (%s) was writted", hookName, hookPath)
		var fileMode os.FileMode = os.ModePerm
		err = c.files.Chmod(hookPath, fileMode)
		if err != nil {
			return err
		}
		log.Debugf("Hook file mode changed to %s", fileMode.String())

		if runtime.GOOS != "windows" {
			err = c.files.Chown(hookPath, c.user)
			if err != nil {
				return err
			}
			log.Debug("Hook file ownership changed to currect user")
		}
	}

	return writeDefaultFishermanConfig(c.files, config.BuildFileConfigPath(c.app.Cwd, c.user, c.mode))
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

func writeDefaultFishermanConfig(accessor i.FileSystem, configPath string) error {
	if !accessor.Exist(configPath) {
		content, err := yaml.Marshal(config.DefaultConfig)
		if err != nil {
			return err
		}

		err = accessor.Write(configPath, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
