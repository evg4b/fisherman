package log_test

import (
	"bytes"
	. "fisherman/pkg/log"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	formattingParams = []interface{}{1, "s", 44.3}
	emptyParamas     = []interface{}{}
)

func TestError(t *testing.T) {
	testData := []struct {
		message string
		output  string
		level   Level
	}{
		{
			message: "error",
			output:  "error\n",
			level:   DebugLevel,
		},
		{
			message: "error",
			output:  "error\n",
			level:   InfoLevel,
		},
		{
			message: "error",
			output:  "error\n",
			level:   ErrorLevel,
		},
		{
			message: "error",
			level:   NoneLevel,
		},
		{
			message: "",
			output:  "\n",
			level:   ErrorLevel,
		},
		{
			message: "\n\n",
			output:  "\n\n\n",
			level:   ErrorLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Error message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output := mockLogModule()
			Configure(OutputConfig{LogLevel: tt.level})
			Error(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestErrorf(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   Level
	}{
		{
			message: "error %d %s %f",
			params:  formattingParams,
			output:  "error 1 s 44.300000\n",
			level:   ErrorLevel,
		},
		{
			message: "error",
			params:  emptyParamas,
			output:  "error\n",
			level:   ErrorLevel,
		},
		{
			message: "error %s",
			params:  emptyParamas,
			output:  "error %!s(MISSING)\n",
			level:   ErrorLevel,
		},
		{
			message: "error %d %s %f",
			params:  formattingParams,
			output:  "error 1 s 44.300000\n",
			level:   DebugLevel,
		},
		{
			message: "error",
			params:  emptyParamas,
			output:  "error\n",
			level:   DebugLevel,
		},
		{
			message: "error %s",
			params:  emptyParamas,
			output:  "error %!s(MISSING)\n",
			level:   DebugLevel,
		},
		{
			message: "error %d %s %f",
			params:  formattingParams,
			output:  "error 1 s 44.300000\n",
			level:   InfoLevel,
		},
		{
			message: "error",
			params:  emptyParamas,
			output:  "error\n",
			level:   InfoLevel,
		},
		{
			message: "error %s",
			params:  emptyParamas,
			output:  "error %!s(MISSING)\n",
			level:   InfoLevel,
		},
		{
			message: "error %d %s %f",
			params:  formattingParams,
			level:   NoneLevel,
		},
		{
			message: "error",
			params:  emptyParamas,
			level:   NoneLevel,
		},
		{
			message: "error %s",
			params:  emptyParamas,
			level:   NoneLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Errorf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})
			Errorf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestDebug(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   Level
	}{
		{
			message: "debug",
			output:  "debug\n",
			level:   DebugLevel,
		},
		{
			message: "debug",
			level:   InfoLevel,
		},
		{
			message: "debug",
			level:   ErrorLevel,
		},
		{
			message: "debug",
			level:   NoneLevel,
		},
		{
			message: "",
			output:  "\n",
			level:   DebugLevel,
		},
		{
			message: "\n\n",
			output:  "\n\n\n",
			level:   DebugLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debug message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})
			Debug(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestDebugf(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   Level
	}{
		{
			message: "debug %d %s %f",
			params:  formattingParams,
			output:  "debug 1 s 44.300000\n",
			level:   DebugLevel,
		},
		{
			message: "debug",
			params:  emptyParamas,
			output:  "debug\n",
			level:   DebugLevel,
		},
		{
			message: "debug %s",
			params:  emptyParamas,
			output:  "debug %!s(MISSING)\n",
			level:   DebugLevel,
		},
		{
			message: "debug %d %s %f",
			params:  formattingParams,
			level:   ErrorLevel,
		},
		{
			message: "debug",
			params:  emptyParamas,
			level:   ErrorLevel,
		},
		{
			message: "debug %s",
			params:  emptyParamas,
			level:   ErrorLevel,
		},
		{
			message: "debug %d %s %f",
			params:  formattingParams,
			level:   InfoLevel,
		},
		{
			message: "debug",
			params:  emptyParamas,
			level:   InfoLevel,
		},
		{
			message: "debug %s",
			params:  emptyParamas,
			level:   InfoLevel,
		},
		{
			message: "debug %d %s %f",
			params:  formattingParams,
			level:   NoneLevel,
		},
		{
			message: "debug",
			params:  emptyParamas,
			level:   NoneLevel,
		},
		{
			message: "debug %s",
			params:  emptyParamas,
			level:   NoneLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debugf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})
			Debugf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestInfo(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   Level
	}{
		{
			message: "info",
			output:  "info\n",
			level:   DebugLevel,
		},
		{
			message: "info",
			output:  "info\n",
			level:   InfoLevel,
		},
		{
			message: "info",
			level:   ErrorLevel,
		},
		{
			message: "info",
			level:   NoneLevel,
		},
		{
			message: "",
			output:  "\n",
			level:   InfoLevel,
		},
		{
			message: "\n\n",
			output:  "\n\n\n",
			level:   InfoLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Info message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})
			Info(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestInfof(t *testing.T) {
	output := mockLogModule()
	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   Level
	}{
		{
			message: "info %d %s %f",
			params:  formattingParams,
			level:   ErrorLevel,
		},
		{
			message: "info",
			params:  emptyParamas,
			level:   ErrorLevel,
		},
		{
			message: "info %s",
			params:  emptyParamas,
			level:   ErrorLevel,
		},
		{
			message: "info %d %s %f",
			params:  formattingParams,
			output:  "info 1 s 44.300000\n",
			level:   DebugLevel,
		},
		{
			message: "info",
			params:  emptyParamas,
			output:  "info\n",
			level:   DebugLevel,
		},
		{
			message: "info %s",
			params:  emptyParamas,
			output:  "info %!s(MISSING)\n",
			level:   DebugLevel,
		},
		{
			message: "info %d %s %f",
			params:  formattingParams,
			output:  "info 1 s 44.300000\n",
			level:   InfoLevel,
		},
		{
			message: "info",
			params:  emptyParamas,
			output:  "info\n",
			level:   InfoLevel,
		},
		{
			message: "info %s",
			params:  emptyParamas,
			output:  "info %!s(MISSING)\n",
			level:   InfoLevel,
		},
		{
			message: "info %d %s %f",
			params:  formattingParams,
			level:   NoneLevel,
		},
		{
			message: "info",
			params:  emptyParamas,
			level:   NoneLevel,
		},
		{
			message: "info %s",
			params:  emptyParamas,
			level:   NoneLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Infof message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})
			Infof(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestStdout(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   Level
	}{
		{
			message: "demo",
			output:  "demo",
			level:   DebugLevel,
		},
		{
			message: "multiline demo\nmultiline demo\nmultiline demo",
			output:  "multiline demo\nmultiline demo\nmultiline demo",
			level:   DebugLevel,
		},
		{
			message: "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			output:  "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			level:   DebugLevel,
		},
		{
			message: "",
			output:  "",
			level:   DebugLevel,
		},
		{
			message: "\t\t\t",
			output:  "\t\t\t",
			level:   DebugLevel,
		},
		{
			message: "test1",
			output:  "test1",
			level:   ErrorLevel,
		},
		{
			message: "test2",
			output:  "test2",
			level:   InfoLevel,
		},
		{
			message: "test3",
			output:  "test3",
			level:   DebugLevel,
		},
		{
			message: "test4",
			output:  "test4",
			level:   NoneLevel,
		},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Write message '%s' correctly", tt.message), func(t *testing.T) {
			output.Reset()
			Configure(OutputConfig{LogLevel: tt.level})

			bytesCount, err := Stdout().Write([]byte(tt.message))

			assert.Equal(t, tt.output, output.String())
			assert.NoError(t, err)
			assert.Equal(t, len([]byte(tt.message)), bytesCount)
		})
	}
}

func mockLogModule() *bytes.Buffer {
	output := bytes.NewBufferString("")
	SetOutput(output)

	return output
}
