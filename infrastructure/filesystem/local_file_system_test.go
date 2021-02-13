package filesystem_test

import (
	"fisherman/constants"
	"fisherman/infrastructure/filesystem"
	"fisherman/testing/testutils"
	"fisherman/utils"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var systemNoFileMessage = ""

func init() {
	if runtime.GOOS == constants.WindowsOS {
		systemNoFileMessage = "The system cannot find the path specified."
	} else {
		systemNoFileMessage = "no such file or directory"
	}
}

func TestLocalFileSystem_Exist(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "exist file", path: existPath, expected: true},
		{name: "not exist file", path: "/test/no/files", expected: false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, fs.Exist(tt.path))
	}
}

func TestLocalFileSystem_Read(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name        string
		path        string
		expected    string
		expectedErr string
	}{
		{
			name:        "exist file",
			path:        existPath,
			expected:    "Hello word",
			expectedErr: "",
		},
		{
			name:        "not exist file",
			path:        "/demo/no/files",
			expected:    "",
			expectedErr: "open /demo/no/files: " + systemNoFileMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := fs.Read(tt.path)

			testutils.CheckError(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLocalFileSystem_Delete(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name        string
		path        string
		expectedErr string
	}{
		{
			name:        "exist file",
			path:        existPath,
			expectedErr: "",
		},
		{
			name:        "not exist file",
			path:        "/demo/no/files",
			expectedErr: "remove /demo/no/files: " + systemNoFileMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Delete(tt.path)

			testutils.CheckError(t, tt.expectedErr, err)
			assert.NoFileExists(t, tt.path)
		})
	}
}

func TestLocalFileSystem_Write(t *testing.T) {
	dir := t.TempDir()
	existPath := writeFile(t, "test.bin", "Some other contents")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "exist file",
			path:     existPath,
			expected: "Some contents",
		},
		{
			name:     "not exist file",
			path:     path.Join(dir, "other-demo.txt"),
			expected: "Some contents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Write(tt.path, tt.expected)

			assert.NoError(t, err)
			assert.FileExists(t, tt.path)
			assert.Equal(t, tt.expected, readFile(t, tt.path))
		})
	}
}

func writeFile(t *testing.T, filepath, content string) string {
	dir := t.TempDir()
	fullPath := path.Join(dir, filepath)
	err := ioutil.WriteFile(fullPath, []byte(content), 0600)
	if err != nil {
		t.Fatal(err)
	}

	return fullPath
}

func readFile(t *testing.T, path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestLocalFileSystem_Reader(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name        string
		path        string
		expected    string
		expectedErr string
	}{
		{
			name:        "exist file",
			path:        existPath,
			expected:    "Hello word",
			expectedErr: "",
		},
		{
			name:        "not exist file",
			path:        "/demo/no/files",
			expected:    "",
			expectedErr: "open /demo/no/files: " + systemNoFileMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := fs.Reader(tt.path)
			if err == nil {
				defer actual.Close()
			}

			testutils.CheckError(t, tt.expectedErr, err)
			if len(tt.expectedErr) == 0 {
				assert.NotNil(t, actual)
				data, err := ioutil.ReadAll(actual)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expected, string(data))
			}
		})
	}
}

func TestLocalFileSystem_Chmod(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	tests := []struct {
		name         string
		path         string
		mode         os.FileMode
		shouldApplay bool
	}{
		{
			name:         "ModePerm",
			path:         existPath,
			mode:         os.ModePerm,
			shouldApplay: true,
		},
		{
			name:         "ModeSetuid",
			path:         existPath,
			mode:         os.ModeSetuid,
			shouldApplay: runtime.GOOS != constants.WindowsOS,
		},
		{
			name:         "ModeSetgid",
			path:         existPath,
			mode:         os.ModeSetgid,
			shouldApplay: runtime.GOOS != constants.WindowsOS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.Chmod(tt.path, tt.mode)

			assert.NoError(t, err)

			info, err := os.Stat(tt.path)
			if err != nil {
				t.Fatal(err)
			}

			isModApplied := (info.Mode() & tt.mode) > 0
			assert.Equal(t, tt.shouldApplay, isModApplied)
		})
	}
}

func TestLocalFileSystem_Chown(t *testing.T) {
	existPath := writeFile(t, "test.txt", "Hello word")

	fs := filesystem.NewLocalFileSystem()

	usr, err := user.Current()
	utils.HandleCriticalError(err)

	assert.NotPanics(t, func() {
		err := fs.Chown(existPath, usr)
		assert.NoError(t, err)
	})
}
