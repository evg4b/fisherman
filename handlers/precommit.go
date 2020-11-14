package handlers

/*
type PreCommitHandler struct{}

func (*PreCommitHandler) IsConfigured(c *config.HooksConfig) bool {
	return c.PreCommitHook.Variables != hooks.Variables{} || len(c.PreCommitHook.Shell) > 0
}

// Handle is a handler for pre-commit hook
func (*PreCommitHandler) Handle(ctx *clicontext.CommandContext, args []string) error {
	config := ctx.Config.PreCommitHook
	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		log.Debugf("Additional variables loading filed: %s", err)

		return err
	}

	config.Compile(ctx.Variables())

	err = shellhandlers.ExecParallel(ctx, ctx.Shell, config.Shell)
	if err != nil {
		return err
	}

	for _, glob := range config.AddFilesToIndex {
		err = ctx.Repository.AddGlob(glob)
		if err != nil {
			return err
		}
	}

	return nil
}
*/
