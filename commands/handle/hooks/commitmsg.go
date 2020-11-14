package hooks

import (
	c "fisherman/config/hooks"
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
	"fisherman/validators"
)

func CommitMsg(factory ctxFactory, conf c.CommitMsgHookConfig, extractor v.VariablesExtractor) *h.HookHandler {
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
