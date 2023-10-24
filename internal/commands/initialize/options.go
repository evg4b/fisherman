package initialize

import (
	"os/user"

	"github.com/go-git/go-billy/v5"
)

type Option = func(c *Command)

func WithCwd(cwd string) Option {
	return func(c *Command) {
		c.cwd = cwd
	}
}

func WithFilesystem(fs billy.Filesystem) Option {
	return func(c *Command) {
		c.fs = fs
	}
}

func WithUser(user *user.User) Option {
	return func(c *Command) {
		c.user = user
	}
}

func WithExecutable(executable string) Option {
	return func(c *Command) {
		c.executable = executable
	}
}
