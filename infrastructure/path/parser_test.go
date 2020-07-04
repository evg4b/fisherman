package path

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRegisteredInPath(t *testing.T) {
	testData := []struct {
		path     string
		app      string
		expected bool
	}{
		{path: makePath("/bin/demo", "/usr/test1"), app: "/bin/demo/fisherman", expected: true},
		{path: makePath("/bin/demo", "/usr/test2", "/bin/demo"), app: "/bin/demo/fisherman.exe", expected: true},
		{path: makePath("/bin/demo", "/usr/test3", "/bin/fisherman"), app: "/bin/demo/demo2/fisherman", expected: false},
		{path: makePath("/dev/demo", "/usr/test4"), app: "/bin/fisherman/fisherman", expected: false},
		{path: makePath("/bin/demo", "/usr/test5"), app: "/bin/fisherman", expected: false},
	}

	for _, tt := range testData {
		t.Run(tt.path, func(t *testing.T) {
			s, err := IsRegisteredInPath(tt.path, tt.app)
			assert.Equal(t, s, tt.expected)
			assert.Nil(t, err)
		})
	}
}

func makePath(paths ...string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}
