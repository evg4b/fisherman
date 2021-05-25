package utils_test

import (
	"fisherman/internal/utils"
	"fisherman/testing/testutils"
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

			actual, err := utils.ReadFileAsString(fs, tt.filepath)

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}
