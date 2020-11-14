package hooks

import (
	"fisherman/actions"
	c "fisherman/config/hooks"
	i "fisherman/infrastructure"
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
)

func PreCommit(factory ctxFactory, conf c.PreCommitHookConfig, extr v.VariablesExtractor, sh i.Shell) *h.HookHandler {
	variables, err := extr.Variables(conf.Variables)
	if err != nil {
		panic(err)
	}

	conf.Compile(variables)

	return h.NewHookHandler(
		factory,
		NoBeforeActions,
		NoSyncValidators,
		scriptWrapper(conf.Shell),
		[]h.Action{
			func(ctx v.SyncValidationContext) (bool, error) {
				return actions.AddToIndex(ctx, conf.AddFilesToIndex)
			},
		},
	)
}
