package initialize

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/pkg/log"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
)

type Command struct {
	flagSet  *flag.FlagSet
	mode     string
	force    bool
	absolute bool
	usage    string
	files    internal.FileSystem
	app      internal.AppInfo
	user     *user.User
}

func NewCommand(files internal.FileSystem, app internal.AppInfo, user *user.User) *Command {
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

func (command *Command) Run(ctx internal.ExecutionContext) error {
	log.Debugf("Statring initialization (force = %t)", command.force)
	if !command.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			exist, err := afero.Exists(command.files, hookPath)
			if err != nil {
				return err
			}

			if exist {
				log.Debugf("Hook '%s' already exist", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		err := result.ErrorOrNil()
		if err != nil {
			return err
		}
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(command.app.Cwd, ".git", "hooks", hookName)

		err := afero.WriteFile(command.files, hookPath, buildHook(hookName, command.getBinaryPath(), command.absolute), os.ModePerm)
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

		err = command.chown(hookPath, command.user)
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

	exist, err := afero.Exists(command.files, configPath)
	if err != nil {
		return err
	}

	if !exist {
		err := afero.WriteFile(command.files, configPath, []byte(configuration.DefaultConfig), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (command *Command) getBinaryPath() string {
	if command.absolute {
		return command.app.Executable
	}

	return constants.AppName
}
