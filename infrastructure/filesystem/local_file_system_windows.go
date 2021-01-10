package filesystem

import (
	"os/user"
)

func (f *LocalFileSystem) Chown(path string, user *user.User) error {
	return nil
}
