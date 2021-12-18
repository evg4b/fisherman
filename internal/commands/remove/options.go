package remove

import (
	"os/user"

	"github.com/go-git/go-billy/v5"
)

type removeOption = func(c *Command)

func WithCwd(cwd string) removeOption {
	return func(c *Command) {
		c.cwd = cwd
	}
}

func WithFilesystem(files billy.Filesystem) removeOption {
	return func(c *Command) {
		c.files = files
	}
}

func WithUser(user *user.User) removeOption {
	return func(c *Command) {
		c.user = user
	}
}

func WithConfigs(configs map[string]string) removeOption {
	return func(c *Command) {
		c.configs = configs
	}
}
