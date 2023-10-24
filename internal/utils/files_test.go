package utils_test

import (
	"fisherman/pkg/guards"
	"fisherman/testing/testutils"
	"path/filepath"
	"testing"

	. "fisherman/internal/utils"

	"github.com/stretchr/testify/assert"
)

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
		name        string
		filepath    string
		expected    bool
		expectedErr string
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
			name:        "not exist file",
			filepath:    ".." + string(filepath.Separator) + "demo",
			expected:    false,
			expectedErr: "chroot boundary crossed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Exists(fs, tt.filepath)

			assert.Equal(t, tt.expected, actual)
			testutils.AssertError(t, tt.expectedErr, err)
		})
	}
}
