package initialize

import (
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal/clicontext"
	"fisherman/utils"
	"flag"
	"fmt"
	"os"
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
}

func NewCommand(handling flag.ErrorHandling) *Command {
	flagSet := flag.NewFlagSet("init", handling)
	command := &Command{
		flagSet: flagSet,
		usage:   "initializes fisherman in git repository",
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

func (c *Command) Run(ctx *clicontext.CommandContext) error {
	log.Debugf("Statring initialization (force = %t)", c.force)
	if !c.force {
		var result *multierror.Error
		for _, hookName := range constants.HooksNames {
			hookPath := filepath.Join(ctx.App.Cwd, ".git", "hooks", hookName)
			log.Debugf("Cheking hook '%s' (%s)", hookName, hookPath)
			if ctx.Files.Exist(hookPath) {
				log.Debugf("Hook '%s' already exist", hookName)
				result = multierror.Append(result, fmt.Errorf("file %s already exists", hookPath))
			}
		}

		return result.ErrorOrNil()
	}

	bin := ctx.App.Executable
	if !c.abslute {
		bin = utils.NormalizePath(bin)
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

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

func writeDefaultFishermanConfig(accessor infrastructure.FileSystem, configPath string) error {
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
