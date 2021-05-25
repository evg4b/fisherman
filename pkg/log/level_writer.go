package log

import (
	"io"
	"io/ioutil"
)

type LevelWriter struct {
	expected Level
	current  Level
	output   io.Writer
}

func NewlevelWriter(output io.Writer, expected Level) *LevelWriter {
	return &LevelWriter{
		expected: expected,
		current:  InfoLevel,
		output:   output,
	}
}

func (wr *LevelWriter) SetLevel(level Level) {
	wr.current = level
}

func (wr *LevelWriter) SetOutput(output io.Writer) {
	wr.output = output
}

func (wr *LevelWriter) Write(p []byte) (n int, err error) {
	if wr.current <= wr.expected {
		return wr.output.Write(p)
	}

	return ioutil.Discard.Write(p)
}
