package prefixwriter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixWriter_Write(t *testing.T) {
	prefix := "prefix: "

	testData := []struct {
		name     string
		input    string
		expected string
	}{
		{"should add prefix for each line", "foo\nbar", "prefix: foo\nprefix: bar"},
		{"should don't add prefix for empty line ", "", ""},
		{"should add prefix for empty line with newline symbol", "\n", "prefix: \n"},
		{"should add prefix for first line", "foo", "prefix: foo"},
		{"should not add prefix for second empty line", "foo\n", "prefix: foo\n"},
	}

	for _, dd := range testData {
		t.Run(dd.name, func(t *testing.T) {
			var buf bytes.Buffer
			prefixwriter := New(&buf, prefix)
			fmt.Fprintf(prefixwriter, dd.input)

			assert.Equal(t, dd.expected, buf.String())
		})
	}
}

func TestPrefixWriter_Write_AppendMode(t *testing.T) {
	prefix := "prefix: "

	var buf bytes.Buffer
	prefixwriter := New(&buf, prefix)

	fmt.Fprintln(prefixwriter, "input 1")
	fmt.Fprint(prefixwriter, "input 2")
	fmt.Fprint(prefixwriter, " with additional string")

	assert.Equal(t, "prefix: input 1\nprefix: input 2 with additional string", buf.String())
}
