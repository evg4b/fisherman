package expression

import (
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defined(t *testing.T) {
	tests := []struct {
		name      string
		arguments []interface{}
		want      interface{}
		wantErr   string
	}{
		{
			name:      "should return error for empty arguments",
			arguments: []interface{}{},
			wantErr:   "incorrect arguments for Defined",
			want:      false,
		},
		{
			name:      "should return false for single nil argument",
			arguments: []interface{}{nil},
			want:      false,
		},
		{
			name:      "should return false for single not nil argument",
			arguments: []interface{}{1},
			want:      true,
		},
		{
			name:      "should return true for multiple not nil arguments",
			arguments: []interface{}{1, 2, 3},
			want:      true,
		},
		{
			name:      "should return true for multiple arguments with nil",
			arguments: []interface{}{1, 2, nil},
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defined(tt.arguments...)
			testutils.CheckError(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
