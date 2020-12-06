package prefixwriter

import (
	"bytes"
	"io"
)

type PrefixWriter struct {
	prefix          string
	writer          io.Writer
	buffer          bytes.Buffer
	trailingNewline bool
}

func New(writer io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{prefix: prefix, writer: writer, trailingNewline: true}
}

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
