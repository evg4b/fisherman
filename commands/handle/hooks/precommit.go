package hooks

import (
	"fisherman/actions"
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
)

func PreCommit(
	ctxFactory internal.CtxFactory,
	configuration hooks.PreCommitHookConfig,
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
		[]handling.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.AddToIndex(ctx, configuration.AddFilesToIndex)
			},
		},
	)

	return HandlerRegistration{Registered: true, Handler: handler}
}
