package initialize

import (
	"os/user"

	"github.com/go-git/go-billy/v5"
)

type initializeOption = func(c *Command)

func WithCwd(cwd string) initializeOption {
	return func(c *Command) {
		c.cwd = cwd
	}
}

func WithFilesystem(fs billy.Filesystem) initializeOption {
	return func(c *Command) {
		c.fs = fs
	}
}

func WithUser(user *user.User) initializeOption {
	return func(c *Command) {
		c.user = user
	}
}

func WithExecutable(executable string) initializeOption {
	return func(c *Command) {
		c.executable = executable
	}
}
