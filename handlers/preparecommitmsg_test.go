package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testBranch = "ref/head/master"

func TestGetPreparedMessage(t *testing.T) {
	testData := []struct {
		message         string
		regexpString    string
		branch          string
		expectedMessage string
		expected        bool
	}{
		{"", "", testBranch, "", false},
		{"MESSAGE", "", testBranch, "MESSAGE", true},
		{"$1", "ref/head/(.*)", testBranch, "master", true},
		{"", "ref/head/(.*)", testBranch, "", false},
		{"$3/$2/$1", "(.*)/(.*)/(.*)", testBranch, "master/head/ref", true},
		{"MESSAGE", "(.*)/(.*)/(.*)", testBranch, "MESSAGE", true},
	}

	for _, tt := range testData {
		t.Run("", func(t *testing.T) {
			message, isPresented := getPreparedMessage(tt.message, tt.regexpString, tt.branch)
			assert.Equal(t, message, tt.expectedMessage)
			assert.Equal(t, isPresented, tt.expected)
		})
	}
}
