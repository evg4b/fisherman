package initialize

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

// Command is structure for storage information about init command
type Command struct {
	flagSet *flag.FlagSet
	mode    string
	force   bool
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Init command created")

	flagSet := flag.NewFlagSet("init", handling)
	command := &Command{flagSet: flagSet}
	modeMessage := fmt.Sprintf("(%s, %s, %s)", config.LocalMode, config.RepoMode, config.GlobalMode)
	flagSet.StringVar(&command.mode, "mode", config.RepoMode, modeMessage)
	flagSet.BoolVar(&command.force, "force", false, "")

	return command
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *clicontext.CommandContext) error {
	log.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if ctx.Files.Exist(hookPath) {
				log.Debugf("Hook '%s' already declared", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		if result.ErrorOrNil() != nil {
			return result
		}
	}

	bin := constants.AppName
	if !ctx.App.IsRegisteredInPath {
		log.Debugf("App is not defined in global scope, will be used '%s' path", ctx.App.Executable)
		bin = fmt.Sprintf("'%s'", ctx.App.Executable)
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
		err := ctx.Files.Write(hookPath, buildHook(bin, hookName))
		if err != nil {
			return err
		}
		log.Infof("Hook '%s' (%s) was writted", hookName, hookPath)
		var fileMode os.FileMode = os.ModePerm
		err = ctx.Files.Chmod(hookPath, fileMode)
		if err != nil {
			return err
		}
		log.Debugf("Hook file mode changed to %s", fileMode.String())

		if runtime.GOOS != "windows" {
			err = ctx.Files.Chown(hookPath, ctx.User)
			if err != nil {
				return err
			}
			log.Debug("Hook file ownership changed to currect user")
		}
	}

	return writeDefaultFishermanConfig(ctx.Files, config.BuildFileConfigPath(ctx.App.Cwd, ctx.User, c.mode))
}

// Name returns command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

func writeDefaultFishermanConfig(accessor infrastructure.FileAccessor, configPath string) error {
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
