package helpers_test

import (
	"fisherman/testing/testutils"
	"testing"

	. "fisherman/pkg/shell/helpers"
)

func TestMergeEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      []string
		vars     map[string]string
		expected []string
	}{
		{
			name:     "empty slice and map",
			env:      []string{},
			vars:     map[string]string{},
			expected: []string{},
		},
		{
			name: "empty source env slice",
			env:  []string{},
			vars: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expected: []string{
				"VAR1=value1",
				"VAR2=value2",
			},
		},
		{
			name: "empty additional vars map",
			env: []string{
				"VAR1=value1",
				"VAR2=value2",
			},
			vars: map[string]string{},
			expected: []string{
				"VAR1=value1",
				"VAR2=value2",
			},
		},
		{
			name: "empty additional vars map",
			env: []string{
				"VAR1=value1",
			},
			vars: map[string]string{
				"VAR2": "value2",
			},
			expected: []string{
				"VAR1=value1",
				"VAR2=value2",
			},
		},
		{
			name: "overwrite source from passed vars",
			env: []string{
				"VAR1=value1",
				"VAR2=value2",
			},
			vars: map[string]string{
				"VAR2": "Custom value",
			},
			expected: []string{
				"VAR1=value1",
				"VAR2=Custom value",
			},
		},
		{
			name: "correctly parse values with =",
			env: []string{
				"VAR1=conation2=test",
				"VAR2=conation2=test2",
				"VAR3=",
			},
			vars: map[string]string{
				"VAR2": "Custom value",
			},
			expected: []string{
				"VAR1=conation2=test",
				"VAR2=Custom value",
				"VAR3=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MergeEnv(tt.env, tt.vars)

			testutils.AssertSlice(t, tt.expected, actual)
		})
	}
}
