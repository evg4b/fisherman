package filesystem

import (
	"os"
	"os/user"
	"strconv"
)

func (f *LocalFileSystem) Chown(path string, user *user.User) error {
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}
