package common

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
)

type CommandExecutionResult struct {
	Key      string
	Output   string
	ExitCode int
	Err      error
}

func ExecCommandsParallel(sh infrastructure.Shell, script hooks.ScriptsConfig) error {
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

	for _, rez := range results {
		log.Debugf("[%s]", rez.Key)
		log.Debug(rez.Output)
	}

	return nil
}

func run(chanel chan CommandExecutionResult, sh infrastructure.Shell, key string, command hooks.ScriptConfig) {
	stdout, exitCode, err := sh.Exec(command.Commands, &command.Env)
	chanel <- CommandExecutionResult{
		Key:      key,
		Output:   stdout,
		Err:      err,
		ExitCode: exitCode,
	}
}
