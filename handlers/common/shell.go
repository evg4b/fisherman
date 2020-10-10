package common

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
)

type CommandExecutionResult struct {
	Key      string
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

type ExecutionPack struct {
	Key    string
	Result CommandExecutionResult
}

func ExecCommandsParallel(sh infrastructure.Shell, commands hooks.CmdConfig) error {
	chanel := make(chan CommandExecutionResult)
	for key, command := range commands {
		log.Debugf("Run cmd %s", key)
		go run(chanel, sh, key, command)
	}

	results := make(map[string]CommandExecutionResult, len(commands))
	for i := 0; i < len(commands); i++ {
		r := <-chanel
		results[r.Key] = r
	}

	return nil
}

func run(chanel chan CommandExecutionResult, sh infrastructure.Shell, key string, command hooks.Command) {
	stdout, stderr, err := sh.Exec(command.Commands, &command.Env)
	chanel <- CommandExecutionResult{
		Key:    key,
		Stdout: stdout,
		Stderr: stderr,
		Err:    err,
	}
}
