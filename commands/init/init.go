package init

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

// Command is structure for storage information about init command
type Command struct {
	fs    *flag.FlagSet
	mode  string
	force bool
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("init", handling)
	c := &Command{fs: fs}
	modeMessage := fmt.Sprintf("(%s, %s, %s)", config.LocalMode, config.RepoMode, config.GlobalMode)
	fs.StringVar(&c.mode, "mode", config.RepoMode, modeMessage)
	fs.BoolVar(&c.force, "force", false, "")
	fs.BoolVar(&c.force, "absolute", false, "")

	return c
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *commands.CommandContext) error {
	logger.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
			logger.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if ctx.Files.Exist(hookPath) {
				logger.Debugf("Hook '%s' already declared", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		if result.ErrorOrNil() != nil {
			return result
		}
	}

	bin := constants.AppName
	if !ctx.App.IsRegisteredInPath {
		logger.Debugf("App is not defined in global scope, will be used '%s' path", ctx.App.Executable)
		bin = fmt.Sprintf("'%s'", ctx.App.Executable)
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
		err := ctx.Files.Write(hookPath, buildHook(bin, hookName))
		utils.HandleCriticalError(err)
		logger.Infof("Hook '%s' (%s) was writted", hookName, hookPath)
	}

	configPath, err := config.BuildFileConfigPath(ctx.App.Cwd, ctx.User, c.mode)
	utils.HandleCriticalError(err)

	err = writeDefaultFishermanConfig(ctx.Files, configPath)
	utils.HandleCriticalError(err)

	return err
}

// Name returns namand name
func (c *Command) Name() string {
	return c.fs.Name()
}

func writeDefaultFishermanConfig(accessor infrastructure.FileAccessor, configPath string) error {
	if !accessor.Exist(configPath) {
		content, err := yaml.Marshal(config.DefaultConfig)
		utils.HandleCriticalError(err)
		err = accessor.Write(configPath, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
