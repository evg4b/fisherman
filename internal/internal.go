package internal

import (
	"context"
	i "fisherman/infrastructure"
	"io"
)

type ExecutionContext interface {
	context.Context
	GlobalVariables() (map[string]interface{}, error)
	Files() i.FileSystem
	Shell() i.Shell
	Repository() i.Repository
	Args() []string
	Output() io.Writer
	Message() (string, error)
	Stop()
}

type AppInfo struct {
	Cwd        string
	Executable string
	Configs    map[string]string
}
