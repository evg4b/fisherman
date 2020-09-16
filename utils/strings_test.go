package utils_test

import (
	"fisherman/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "Return true for empty string", value: "", want: true},
		{name: "Return true for spaces only", value: "  ", want: true},
		{name: "Return true for tabs only", value: "\t\t", want: true},
		{name: "Return true for newline charsets only", value: "\n\n", want: true},
		{name: "Return true for mixed witespace string", value: "\t\n ", want: true},
		{name: "Return false for strign without space", value: "2", want: false},
		{name: "Return true for not empty string with space", value: " 2", want: false},
		{name: "Return true for not empty string with tab", value: "\t2", want: false},
		{name: "Return true for not empty string with newline", value: "\n2", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.IsEmpty(tt.value))
		})
	}
}
