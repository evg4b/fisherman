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
) *handling.HookHandler {
	compile(&configuration)

	return handling.NewHookHandler(
		ctxFactory,
		NoBeforeActions,
		NoSyncValidators,
		scriptWrapper(configuration.Shell),
		NoAfterActions,
	)
}
