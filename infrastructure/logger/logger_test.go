package logger_test

import (
	"bytes"
	"fisherman/infrastructure/logger"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var formattingParams = []interface{}{1, "s", 44.3}
var emptyParamas = []interface{}{}

func TestErrorWithLogLevel(t *testing.T) {
	t.Parallel()

	testData := []struct {
		message string
		output  string
		level   logger.LogLevel
	}{
		{message: "error", output: "error\n", level: logger.DebugLevel},
		{message: "error", output: "error\n", level: logger.InfoLevel},
		{message: "error", output: "error\n", level: logger.ErrorLevel},
		{message: "error", output: "", level: logger.NoneLevel},
		{message: "", output: "\n", level: logger.ErrorLevel},
		{message: "\n\n", output: "\n\n\n", level: logger.ErrorLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Error message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output := mockLogModule()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Error(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestErrorfWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   logger.LogLevel
	}{
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: logger.ErrorLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: logger.ErrorLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: logger.ErrorLevel},
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: logger.DebugLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: logger.DebugLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: logger.DebugLevel},
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: logger.InfoLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: logger.InfoLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: logger.InfoLevel},
		{message: "error %d %s %f", params: formattingParams, output: "", level: logger.NoneLevel},
		{message: "error", params: emptyParamas, output: "", level: logger.NoneLevel},
		{message: "error %s", params: emptyParamas, output: "", level: logger.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Errorf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Errorf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestDebugWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   logger.LogLevel
	}{
		{message: "debug", output: "debug\n", level: logger.DebugLevel},
		{message: "debug", output: "", level: logger.InfoLevel},
		{message: "debug", output: "", level: logger.ErrorLevel},
		{message: "debug", output: "", level: logger.NoneLevel},
		{message: "", output: "\n", level: logger.DebugLevel},
		{message: "\n\n", output: "\n\n\n", level: logger.DebugLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debug message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Debug(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestDebugfWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   logger.LogLevel
	}{
		{message: "debug %d %s %f", params: formattingParams, output: "debug 1 s 44.300000\n", level: logger.DebugLevel},
		{message: "debug", params: emptyParamas, output: "debug\n", level: logger.DebugLevel},
		{message: "debug %s", params: emptyParamas, output: "debug %!s(MISSING)\n", level: logger.DebugLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: logger.ErrorLevel},
		{message: "debug", params: emptyParamas, output: "", level: logger.ErrorLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: logger.ErrorLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: logger.InfoLevel},
		{message: "debug", params: emptyParamas, output: "", level: logger.InfoLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: logger.InfoLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: logger.NoneLevel},
		{message: "debug", params: emptyParamas, output: "", level: logger.NoneLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: logger.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debugf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Debugf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestInfoWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   logger.LogLevel
	}{
		{message: "info", output: "info\n", level: logger.DebugLevel},
		{message: "info", output: "info\n", level: logger.InfoLevel},
		{message: "info", output: "", level: logger.ErrorLevel},
		{message: "info", output: "", level: logger.NoneLevel},
		{message: "", output: "\n", level: logger.InfoLevel},
		{message: "\n\n", output: "\n\n\n", level: logger.InfoLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Info message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Info(tt.message)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestInfofWithLogLevel(t *testing.T) {
	output := mockLogModule()
	testData := []struct {
		message string
		output  string
		params  []interface{}
		level   logger.LogLevel
	}{
		{message: "info %d %s %f", params: formattingParams, output: "", level: logger.ErrorLevel},
		{message: "info", params: emptyParamas, output: "", level: logger.ErrorLevel},
		{message: "info %s", params: emptyParamas, output: "", level: logger.ErrorLevel},
		{message: "info %d %s %f", params: formattingParams, output: "info 1 s 44.300000\n", level: logger.DebugLevel},
		{message: "info", params: emptyParamas, output: "info\n", level: logger.DebugLevel},
		{message: "info %s", params: emptyParamas, output: "info %!s(MISSING)\n", level: logger.DebugLevel},
		{message: "info %d %s %f", params: formattingParams, output: "info 1 s 44.300000\n", level: logger.InfoLevel},
		{message: "info", params: emptyParamas, output: "info\n", level: logger.InfoLevel},
		{message: "info %s", params: emptyParamas, output: "info %!s(MISSING)\n", level: logger.InfoLevel},
		{message: "info %d %s %f", params: formattingParams, output: "", level: logger.NoneLevel},
		{message: "info", params: emptyParamas, output: "", level: logger.NoneLevel},
		{message: "info %s", params: emptyParamas, output: "", level: logger.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Infof message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: tt.level})
			logger.Infof(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestWrite(t *testing.T) {
	output := mockLogModule()

	testData := []string{
		"demo",
		"multiline demo\nmultiline demo\nmultiline demo",
		"multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
		"",
		"\t\t\t",
	}

	for _, message := range testData {
		t.Run(fmt.Sprintf("Write message '%s' correctly", message), func(t *testing.T) {
			output.Reset()
			logger.Configure(logger.OutputConfig{LogLevel: logger.DebugLevel})
			bytesCount, err := logger.Writer().Write([]byte(message))
			assert.Equal(t, message, output.String())
			assert.NoError(t, err)
			assert.Equal(t, len([]byte(message)), bytesCount)
		})
	}
}

func mockLogModule() *bytes.Buffer {
	output := bytes.NewBufferString("")
	log.SetOutput(output)
	log.SetFlags(0)

	return output
}
