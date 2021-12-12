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

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Command) Run(ctx context.Context) error {
	handler, err := handling.NewHookHandler(
		c.hook,
		handling.WithExpressionEngine(c.engine),
		handling.WithHooksConfig(c.config),
		handling.WithGlobalVars(c.globalVars),
		handling.WithCwd(c.cwd),
		handling.WithFileSystem(c.fs),
		handling.WithRepository(c.repo),
		handling.WithArgs(c.args),
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

	utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
		constants.HookName:                 c.hook,
		constants.GlobalConfigPath:         utils.OriginalOrNA(c.configFiles[configuration.GlobalMode]),
		constants.RepoConfigPath:           utils.OriginalOrNA(c.configFiles[configuration.RepoMode]),
		constants.LocalConfigPath:          utils.OriginalOrNA(c.configFiles[configuration.LocalMode]),
		constants.FishermanVersionVariable: constants.Version,
	})

	return handler.Handle(ctx)
}
