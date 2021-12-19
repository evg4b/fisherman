package remove

import (
	"github.com/go-git/go-billy/v5"
)

type removeOption = func(c *Command)

func WithCwd(cwd string) removeOption {
	return func(c *Command) {
		c.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) removeOption {
	return func(c *Command) {
		c.fs = fs
	}
}

func WithConfigFiles(configs map[string]string) removeOption {
	return func(c *Command) {
		c.configFiles = configs
	}
}
