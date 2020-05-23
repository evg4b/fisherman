package path

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsRegisteredInPath(t *testing.T) {
	testData := []struct {
		path     string
		app      string
		expected bool
	}{
		{path: "/bin/demo;/usr/test1;", app: "/bin/demo/fisherman", expected: true},
		{path: "/bin/demo;/usr/test2;", app: "/bin/demo/fisherman.exe", expected: true},
		{path: "/bin/demo;/usr/test3;", app: "/bin/demo/demo2/fisherman", expected: false},
		{path: "/dev/demo;/usr/test4;", app: "/bin/fisherman/fisherman", expected: false},
		{path: "/bin/demo;/usr/test5;", app: "/bin/fisherman", expected: false},
	}

	for _, tt := range testData {
		t.Run(tt.path, func(t *testing.T) {
			s, err := IsRegisteredInPath(tt.path, tt.app)
			assert.Equal(t, s, tt.expected)
			assert.Nil(t, err)
		})
	}
}
