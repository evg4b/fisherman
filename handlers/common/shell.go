package common

import (
	"fisherman/config/hooks"
	inf "fisherman/infrastructure"
	"fisherman/infrastructure/log"
)

type CommandExecutionResult struct {
	Key      string
	Output   string
	ExitCode int
	Err      error
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
	stdout, exitCode, err := sh.Exec(command.Commands, &command.Env, command.Output)
	chanel <- CommandExecutionResult{
		Key:      key,
		Output:   stdout,
		Err:      err,
		ExitCode: exitCode,
	}
}
