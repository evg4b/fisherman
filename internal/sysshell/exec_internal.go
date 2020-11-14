package sysshell

import (
	"context"
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/shell"
	"sync"
)

func execInternal(chanel chan shell.ExecResult, ctx ContextWithStop, sh inf.Shell, scripts hooks.ScriptsConfig) {
	var wg sync.WaitGroup
	for scriptName, shellScript := range scripts {
		wg.Add(1)
		go func(name string, script hooks.ScriptConfig) {
			defer wg.Done()

			result := sh.Exec(ctx, shell.ScriptConfig{
				Name:     name,
				Commands: script.Commands,
				Env:      script.Env,
				Output:   script.Output,
			})

			if context.Canceled == ctx.Err() {
				result.Error = ctx.Err()
			}

			chanel <- result

			if !result.IsSuccessful() {
				ctx.Stop()
			}
		}(scriptName, shellScript)
	}
	wg.Wait()
	close(chanel)
}
