package prefixwriter

import (
	"bytes"
	"io"
)

// PrefixWriter is io.Writer interface implementation where the specified prefix
// will be added to the beginning of a new line.
type PrefixWriter struct {
	prefix          string
	writer          io.Writer
	buffer          bytes.Buffer
	trailingNewline bool
}

// NewWriter wrappers passed io.Writer to PrefixWriter with passed prefix.
func NewWriter(writer io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{prefix: prefix, writer: writer, trailingNewline: true}
}

// Write is io.Writer interface implementation.
//
// Write writes len(payload) bytes from payload to the underlying data stream.
// It returns the number of bytes written from payload (0 <= n <= len(payload))
// and any error encountered that caused the write to stop early.
// Write returns a non-nil error if it returns n < len(payload).
// Write does not modify the slice data, even temporarily.
// Write does not retain payload.
func (writer *PrefixWriter) Write(payload []byte) (int, error) {
	writer.buffer.Reset()

	for _, b := range payload {
		if writer.trailingNewline {
			writer.buffer.WriteString(writer.prefix)
			writer.trailingNewline = false
		}

		writer.buffer.WriteByte(b)

		if b == '\n' {
			writer.trailingNewline = true
		}
	}

	_, err := writer.writer.Write(writer.buffer.Bytes())

	return len(payload), err
}
