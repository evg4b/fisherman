package init

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/utils"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

// Command is structure for storage information about init command
type Command struct {
	fs       *flag.FlagSet
	mode     string
	absolute bool
	force    bool
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

// Run executes init command
func (c *Command) Run(ctx *context.CommandContext, args []string) error {
	c.fs.Parse(args)
	ctx.Logger.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(ctx.AppInfo.Cwd, ".git", "hooks", hookName)
			ctx.Logger.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if ctx.FileAccessor.FileExist(hookPath) {
				ctx.Logger.Debugf("Hook '%s' already declared", hookName)
				result = multierror.Append(result, fmt.Errorf("File %s already exists", hookPath))
			}
		}

		if result.ErrorOrNil() != nil {
			return result
		}
	}

	bin := constants.AppName
	if !ctx.AppInfo.IsRegisteredInPath {
		ctx.Logger.Debugf("App is not defined in global scope, will be used '%s' path", ctx.AppInfo.AppPath)
		bin = ctx.AppInfo.AppPath
	}

	for _, hookName := range constants.HooksNames {
		hookPath := filepath.Join(ctx.AppInfo.Cwd, ".git", "hooks", hookName)
		err := ctx.FileAccessor.WriteFile(hookPath, buildHook(bin, hookName))
		utils.HandleCriticalError(err)
		ctx.Logger.Debugf("Hook '%s' (%s) was writted", hookName, hookPath)
	}

	configPath, err := config.BuildFileConfigPath(ctx.AppInfo.Cwd, ctx.User, c.mode)
	utils.HandleCriticalError(err)

	err = writeFishermanConfig(ctx.FileAccessor, configPath)
	utils.HandleCriticalError(err)

	return err
}

// Name returns namand name
func (c *Command) Name() string {
	return c.fs.Name()
}
