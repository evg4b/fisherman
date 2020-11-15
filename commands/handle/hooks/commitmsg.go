package hooks

import (
	"fisherman/config/hooks"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/validators"
)

func CommitMsg(
	ctxFactory internal.CtxFactory,
	configuration hooks.CommitMsgHookConfig,
	compile configcompiler.Compiler,
) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotResigter
	}

	compile(&configuration)

	handler := handling.NewHookHandler(
		ctxFactory,
		NoBeforeActions,
		[]validation.SyncValidator{
			boolWrapper(validators.MessageNotEmpty, configuration.NotEmpty),
			stringWrapper(validators.MessageHasPrefix, configuration.MessagePrefix),
			stringWrapper(validators.MessageHasSuffix, configuration.MessageSuffix),
			stringWrapper(validators.MessageRegexp, configuration.MessageRegexp),
		},
		NoAsyncValidators,
		NoAfterActions,
	)

	return HandlerRegistration{Registered: true, Handler: handler}
}
