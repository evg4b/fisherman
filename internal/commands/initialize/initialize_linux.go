package initialize

import (
	"os/user"
	"strconv"
)

func (command *Command) chown(path string, user *user.User) error {
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}

	return command.files.Chown(path, uid, gid)
}
