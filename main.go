package main

import (
	"context"
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/internal/app"
	"github.com/evg4b/fisherman/internal/commands/handle"
	"github.com/evg4b/fisherman/internal/commands/initialize"
	"github.com/evg4b/fisherman/internal/commands/remove"
	"github.com/evg4b/fisherman/internal/commands/version"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/expression"
	"github.com/evg4b/fisherman/internal/utils"
	"github.com/evg4b/fisherman/pkg/guards"
	"github.com/evg4b/fisherman/pkg/log"
	"github.com/evg4b/fisherman/pkg/vcs"
	"os"
	"os/user"
	"runtime"

	"github.com/go-git/go-billy/v5/osfs"
)

const fatalExitCode = 1

func main() {
	defer utils.PanicInterceptor(func(recovered any) {
		log.Errorf("Fatal error: %s", recovered)
		if err, ok := recovered.(error); ok {
			log.DumpError(err)
		}

		os.Exit(fatalExitCode)
	})

	usr, err := user.Current()
	guards.NoError(err)

	cwd, err := os.Getwd()
	guards.NoError(err)

	executablePath, err := os.Executable()
	guards.NoError(err)

	fs := osfs.New("")

	configs, err := configuration.FindConfigFiles(usr, cwd, fs)
	guards.NoError(err)

	config, err := configuration.Load(fs, configs)
	guards.NoError(err)

	log.Configure(config.Output)

	ctx := context.Background()
	engine := expression.NewGoExpressionEngine()
	repo := vcs.NewRepository(vcs.WithFsRepo(cwd))

	fishermanApp := app.NewFishermanApp(
		app.WithCommands([]internal.CliCommand{
			initialize.NewCommand(
				initialize.WithCwd(cwd),
				initialize.WithFilesystem(fs),
				initialize.WithUser(usr),
				initialize.WithExecutable(executablePath),
			),
			handle.NewCommand(
				handle.WithExpressionEngine(engine),
				handle.WithHooksConfig(&config.Hooks),
				handle.WithGlobalVars(config.GlobalVariables),
				handle.WithCwd(cwd),
				handle.WithFileSystem(fs),
				handle.WithRepository(repo),
				handle.WithEnv(os.Environ()),
				handle.WithWorkersCount(uint(runtime.NumCPU())),
				handle.WithConfigFiles(configs),
				handle.WithOutput(os.Stdout),
			),
			remove.NewCommand(
				remove.WithCwd(cwd),
				remove.WithFileSystem(fs),
				remove.WithConfigFiles(configs),
			),
			version.NewCommand(),
		}),
		app.WithSistemInterruptSignals(),
	)

	if err = fishermanApp.Run(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}
