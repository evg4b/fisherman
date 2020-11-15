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
) *handling.HookHandler {
	compile(&configuration)

	return handling.NewHookHandler(
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
}
