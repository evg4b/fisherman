package testutils

import (
	"os/user"
	"path/filepath"
)

var TestUser = user.User{
	Uid:      "1",
	Gid:      "2",
	Username: "evg4b",
	Name:     "Evgeny Abramovitch",
	HomeDir:  filepath.Join("usr", "home"),
}
