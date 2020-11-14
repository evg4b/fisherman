package hooks

import (
	c "fisherman/config/hooks"
	"fisherman/internal"
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
	"fisherman/validators"
)

func CommitMsg(factory internal.CtxFactory, conf c.CommitMsgHookConfig, extractor v.VarExtractor) *h.HookHandler {
	variables, err := extractor.Variables(conf.Variables)
	if err != nil {
		panic(err)
	}

	conf.Compile(variables)

	return h.NewHookHandler(
		factory,
		NoBeforeActions,
		[]v.SyncValidator{
			boolWrapper(validators.MessageNotEmpty, conf.NotEmpty),
			stringWrapper(validators.MessageHasPrefix, conf.MessagePrefix),
			stringWrapper(validators.MessageHasSuffix, conf.MessageSuffix),
			stringWrapper(validators.MessageRegexp, conf.MessageRegexp),
		},
		NoAsyncValidators,
		NoAfterActions,
	)
}
