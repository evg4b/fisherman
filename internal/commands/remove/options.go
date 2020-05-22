package remove

import (
	"github.com/go-git/go-billy/v5"
)

type RemoveOption = func(c *Command)

func WithCwd(cwd string) RemoveOption {
	return func(c *Command) {
		c.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) RemoveOption {
	return func(c *Command) {
		c.fs = fs
	}
}

func WithConfigFiles(configs map[string]string) RemoveOption {
	return func(c *Command) {
		c.configFiles = configs
	}
}
