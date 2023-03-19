package handle

import (
	"context"
	"errors"
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/handling"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
)

const noFilesLabel = "N/A"

func (c *Command) Run(ctx context.Context, args []string) error {
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	handler, err := handling.NewHookHandler(
		c.hook,
		handling.WithExpressionEngine(c.engine),
		handling.WithHooksConfig(c.config),
		handling.WithGlobalVars(c.globalVars),
		handling.WithCwd(c.cwd),
		handling.WithFileSystem(c.fs),
		handling.WithRepository(c.repo),
		handling.WithArgs(c.flagSet.Args()),
		handling.WithEnv(c.env),
		handling.WithWorkersCount(c.workersCount),
		handling.WithOutput(c.output),
	)
	if err != nil {
		if errors.Is(err, handling.ErrNotPresented) {
			log.Debugf("hook %s not presented", c.hook)

			return nil
		}

		return err
	}

	log.Debugf("handler for '%s' hook founded", c.hook)

	files := c.configFiles
	utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
		constants.HookName:                 c.hook,
		constants.GlobalConfigPath:         utils.FirstNotEmpty(files[configuration.GlobalMode], noFilesLabel),
		constants.RepoConfigPath:           utils.FirstNotEmpty(files[configuration.RepoMode], noFilesLabel),
		constants.LocalConfigPath:          utils.FirstNotEmpty(files[configuration.LocalMode], noFilesLabel),
		constants.FishermanVersionVariable: constants.Version,
	})

	return handler.Handle(ctx)
}
