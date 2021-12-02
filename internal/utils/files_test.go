package utils_test

import (
	"errors"
	. "fisherman/internal/utils"
	"fisherman/pkg/guards"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileAsString(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		filepath      string
		expected      string
		expectedError string
	}{
		{
			name:          "there are no files",
			files:         map[string]string{},
			filepath:      "demo.txt",
			expected:      "",
			expectedError: "open demo.txt: file does not exist",
		},
		{
			name:          "file not exists",
			files:         map[string]string{"demo2.txt": "content"},
			filepath:      "demo.txt",
			expected:      "",
			expectedError: "open demo.txt: file does not exist",
		},
		{
			name:          "file exists",
			files:         map[string]string{"demo.txt": "content"},
			filepath:      "demo.txt",
			expected:      "content",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := testutils.FsFromMap(t, tt.files)

			actual, err := ReadFileAsString(fs, tt.filepath)

			assert.Equal(t, tt.expected, actual)
			testutils.AssertError(t, tt.expectedError, err)
		})
	}
}

func TestExists(t *testing.T) {
	fs := testutils.FsFromMap(t, map[string]string{
		"demo.txt":   "content",
		"locked.txt": "locked file",
	})

	f, err := fs.Open("locked.txt")
	guards.NoError(err)
	defer func() {
		err := f.Close()
		guards.NoError(err)
	}()

	err = f.Lock()
	guards.NoError(err)

	defer func() {
		err := f.Unlock()
		guards.NoError(err)
	}()

	tests := []struct {
		name          string
		filepath      string
		expected      bool
		expectedError string
	}{
		{
			name:     "file exists",
			filepath: "demo.txt",
			expected: true,
		},
		{
			name:     "locked file",
			filepath: "locked.txt",
			expected: true,
		},
		{
			name:     "not exist file",
			filepath: "notexist.txt",
			expected: false,
		},
		{
			name:          "not exist file",
			filepath:      ".." + string(filepath.Separator) + "demo",
			expected:      false,
			expectedError: "chroot boundary crossed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Exists(fs, tt.filepath)

			assert.Equal(t, tt.expected, actual)
			testutils.AssertError(t, tt.expectedError, err)
		})
	}
}

func TestReadingFailed(t *testing.T) {
	fileMock := mocks.NewFileMock(t).
		ReadMock.Return(0, errors.New("test error"))
	fs := mocks.NewFilesystemMock(t).
		OpenMock.Return(fileMock, nil)

	content, err := ReadFileAsString(fs, "demo.txt")

	assert.Empty(t, content)
	assert.EqualError(t, err, "test error")
}
