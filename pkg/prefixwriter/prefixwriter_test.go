package prefixwriter_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/evg4b/fisherman/pkg/prefixwriter"

	"github.com/stretchr/testify/assert"
)

func TestPrefixWriter_Write(t *testing.T) {
	prefix := "prefix: "

	t.Run("write full line", func(t *testing.T) {
		testData := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "should add prefix for each line",
				input:    "foo\nbar",
				expected: "prefix: foo\nprefix: bar",
			},
			{
				name:     "should don't add prefix for empty line ",
				input:    "",
				expected: "",
			},
			{
				name:     "should add prefix for empty line with newline symbol",
				input:    "\n",
				expected: "prefix: \n",
			},
			{
				name:     "should add prefix for first line",
				input:    "foo",
				expected: "prefix: foo",
			},
			{
				name:     "should not add prefix for second empty line",
				input:    "foo\n",
				expected: "prefix: foo\n",
			},
		}

		for _, testCase := range testData {
			testCase := testCase
			t.Run(testCase.name, func(t *testing.T) {
				var buf bytes.Buffer
				fmt.Fprintf(NewWriter(&buf, prefix), testCase.input)

				assert.Equal(t, testCase.expected, buf.String())
			})
		}
	})

	t.Run("write partial", func(t *testing.T) {
		var buf bytes.Buffer
		prefixWriter := NewWriter(&buf, prefix)

		fmt.Fprintln(prefixWriter, "input 1")
		fmt.Fprint(prefixWriter, "input 2")
		fmt.Fprint(prefixWriter, " with additional string")

		assert.Equal(t, "prefix: input 1\nprefix: input 2 with additional string", buf.String())
	})
}
