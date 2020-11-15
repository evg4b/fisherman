package hooks

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
)

func PrePush(
	ctxFactory internal.CtxFactory,
	configuration hooks.PrePushHookConfig,
	sysShell infrastructure.Shell,
	compile configcompiler.Compiler,
) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotResigter
	}

	compile(&configuration)

	handler := handling.NewHookHandler(
		ctxFactory,
		NoBeforeActions,
		NoSyncValidators,
		scriptWrapper(configuration.Shell),
		NoAfterActions,
	)

	return HandlerRegistration{Registered: true, Handler: handler}
}
