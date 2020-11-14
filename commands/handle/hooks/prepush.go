package hooks

import (
	c "fisherman/config/hooks"
	i "fisherman/infrastructure"
	"fisherman/internal"
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
)

func PrePush(factory internal.CtxFactory, conf c.PrePushHookConfig, extr v.VarExtractor, sh i.Shell) *h.HookHandler {
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
		NoAfterActions,
	)
}
