package common

import (
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/infrastructure/shell"
)

type CommandExecutionResult struct {
	Key    string
	Result shell.ExecResult
}

func ExecCommandsParallel(sh inf.Shell, script hooks.ScriptsConfig) map[string]CommandExecutionResult {
	chanel := make(chan CommandExecutionResult)
	for key, command := range script {
		log.Debugf("Run cmd %s", key)
		go run(chanel, sh, key, command)
	}

	results := make(map[string]CommandExecutionResult, len(script))
	for i := 0; i < len(script); i++ {
		r := <-chanel
		results[r.Key] = r
	}

	return results
}

func run(chanel chan CommandExecutionResult, sh inf.Shell, key string, command hooks.ScriptConfig) {
	result := sh.Exec(command.Commands, &command.Env, command.Output)
	chanel <- CommandExecutionResult{
		Key:    key,
		Result: result,
	}
}
