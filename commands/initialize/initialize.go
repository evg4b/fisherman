package initialize

import (
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/utils"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type Command struct {
	flagSet  *flag.FlagSet
	mode     string
	force    bool
	absolute bool
	usage    string
	files    infrastructure.FileSystem
	app      *internal.AppInfo
	user     *user.User
}

func NewCommand(files infrastructure.FileSystem, app *internal.AppInfo, user *user.User) *Command {
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

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run() error {
	log.Debugf("Statring initialization (force = %t)", command.force)
	if !command.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if command.files.Exist(hookPath) {
				log.Debugf("Hook '%s' already exist", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		err := result.ErrorOrNil()
		if err != nil {
			return err
		}
	}

	bin := command.app.Executable
	if !command.absolute {
		bin = utils.NormalizePath(bin)
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)
		err := command.files.Write(hookPath, buildHook(bin, hookName))
		if err != nil {
			return err
		}
		log.Infof("Hook '%s' (%s) was writted", hookName, hookPath)
		var fileMode os.FileMode = os.ModePerm
		err = command.files.Chmod(hookPath, fileMode)
		if err != nil {
			return err
		}
		log.Debugf("Hook file mode changed to %s", fileMode.String())

		err = command.files.Chown(hookPath, command.user)
		if err != nil {
			return err
		}
		log.Debug("Hook file ownership changed to currect user")
	}

	return command.writeConfig()
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}

func (command *Command) writeConfig() error {
	configFolder := configuration.GetConfigFolder(command.user, command.app.Cwd, command.mode)
	configPath := filepath.Join(configFolder, constants.AppConfigNames[0])

	if !command.files.Exist(configPath) {
		content, err := yaml.Marshal(configuration.DefaultConfig)
		if err != nil {
			return err
		}

		err = command.files.Write(configPath, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
