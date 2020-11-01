package log_test

//nolint: depguard
import (
	"bytes"
	"fisherman/infrastructure/log"
	"fmt"

	syslog "log"
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
		level   log.Level
	}{
		{message: "error", output: "error\n", level: log.DebugLevel},
		{message: "error", output: "error\n", level: log.InfoLevel},
		{message: "error", output: "error\n", level: log.ErrorLevel},
		{message: "error", output: "", level: log.NoneLevel},
		{message: "", output: "\n", level: log.ErrorLevel},
		{message: "\n\n", output: "\n\n\n", level: log.ErrorLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Error message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output := mockLogModule()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Error(tt.message)
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
		level   log.Level
	}{
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: log.ErrorLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: log.ErrorLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: log.ErrorLevel},
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: log.DebugLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: log.DebugLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: log.DebugLevel},
		{message: "error %d %s %f", params: formattingParams, output: "error 1 s 44.300000\n", level: log.InfoLevel},
		{message: "error", params: emptyParamas, output: "error\n", level: log.InfoLevel},
		{message: "error %s", params: emptyParamas, output: "error %!s(MISSING)\n", level: log.InfoLevel},
		{message: "error %d %s %f", params: formattingParams, output: "", level: log.NoneLevel},
		{message: "error", params: emptyParamas, output: "", level: log.NoneLevel},
		{message: "error %s", params: emptyParamas, output: "", level: log.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Errorf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Errorf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestDebugWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   log.Level
	}{
		{message: "debug", output: "debug\n", level: log.DebugLevel},
		{message: "debug", output: "", level: log.InfoLevel},
		{message: "debug", output: "", level: log.ErrorLevel},
		{message: "debug", output: "", level: log.NoneLevel},
		{message: "", output: "\n", level: log.DebugLevel},
		{message: "\n\n", output: "\n\n\n", level: log.DebugLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debug message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Debug(tt.message)
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
		level   log.Level
	}{
		{message: "debug %d %s %f", params: formattingParams, output: "debug 1 s 44.300000\n", level: log.DebugLevel},
		{message: "debug", params: emptyParamas, output: "debug\n", level: log.DebugLevel},
		{message: "debug %s", params: emptyParamas, output: "debug %!s(MISSING)\n", level: log.DebugLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: log.ErrorLevel},
		{message: "debug", params: emptyParamas, output: "", level: log.ErrorLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: log.ErrorLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: log.InfoLevel},
		{message: "debug", params: emptyParamas, output: "", level: log.InfoLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: log.InfoLevel},
		{message: "debug %d %s %f", params: formattingParams, output: "", level: log.NoneLevel},
		{message: "debug", params: emptyParamas, output: "", level: log.NoneLevel},
		{message: "debug %s", params: emptyParamas, output: "", level: log.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Debugf message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Debugf(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

func TestInfoWithLogLevel(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   log.Level
	}{
		{message: "info", output: "info\n", level: log.DebugLevel},
		{message: "info", output: "info\n", level: log.InfoLevel},
		{message: "info", output: "", level: log.ErrorLevel},
		{message: "info", output: "", level: log.NoneLevel},
		{message: "", output: "\n", level: log.InfoLevel},
		{message: "\n\n", output: "\n\n\n", level: log.InfoLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Info message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Info(tt.message)
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
		level   log.Level
	}{
		{message: "info %d %s %f", params: formattingParams, output: "", level: log.ErrorLevel},
		{message: "info", params: emptyParamas, output: "", level: log.ErrorLevel},
		{message: "info %s", params: emptyParamas, output: "", level: log.ErrorLevel},
		{message: "info %d %s %f", params: formattingParams, output: "info 1 s 44.300000\n", level: log.DebugLevel},
		{message: "info", params: emptyParamas, output: "info\n", level: log.DebugLevel},
		{message: "info %s", params: emptyParamas, output: "info %!s(MISSING)\n", level: log.DebugLevel},
		{message: "info %d %s %f", params: formattingParams, output: "info 1 s 44.300000\n", level: log.InfoLevel},
		{message: "info", params: emptyParamas, output: "info\n", level: log.InfoLevel},
		{message: "info %s", params: emptyParamas, output: "info %!s(MISSING)\n", level: log.InfoLevel},
		{message: "info %d %s %f", params: formattingParams, output: "", level: log.NoneLevel},
		{message: "info", params: emptyParamas, output: "", level: log.NoneLevel},
		{message: "info %s", params: emptyParamas, output: "", level: log.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Infof message '%s' for level: %d", tt.message, tt.level), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			log.Infof(tt.message, tt.params...)
			assert.Equal(t, tt.output, output.String())
		})
	}
}

//nolint:dupl
func TestWrite(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   log.Level
	}{
		{message: "demo", output: "demo", level: log.DebugLevel},
		{
			message: "multiline demo\nmultiline demo\nmultiline demo",
			output:  "multiline demo\nmultiline demo\nmultiline demo",
			level:   log.DebugLevel,
		},
		{
			message: "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			output:  "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			level:   log.DebugLevel,
		},
		{message: "", output: "", level: log.DebugLevel},
		{message: "\t\t\t", output: "\t\t\t", level: log.DebugLevel},
		{message: "test1", output: "", level: log.ErrorLevel},
		{message: "test2", output: "test2", level: log.InfoLevel},
		{message: "test3", output: "test3", level: log.DebugLevel},
		{message: "test4", output: "", level: log.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Write message '%s' correctly", tt.message), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			bytesCount, err := log.Writer().Write([]byte(tt.message))
			assert.Equal(t, tt.output, output.String())
			assert.NoError(t, err)
			assert.Equal(t, len([]byte(tt.message)), bytesCount)
		})
	}
}

//nolint:dupl
func TestRawWriter(t *testing.T) {
	output := mockLogModule()

	testData := []struct {
		message string
		output  string
		level   log.Level
	}{
		{message: "demo", output: "demo", level: log.DebugLevel},
		{
			message: "multiline demo\nmultiline demo\nmultiline demo",
			output:  "multiline demo\nmultiline demo\nmultiline demo",
			level:   log.DebugLevel,
		},
		{
			message: "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			output:  "multiline demo win\r\nmultiline demo win\r\nmultiline demo win",
			level:   log.DebugLevel,
		},
		{message: "", output: "", level: log.DebugLevel},
		{message: "\t\t\t", output: "\t\t\t", level: log.DebugLevel},
		{message: "test1", output: "test1", level: log.ErrorLevel},
		{message: "test2", output: "test2", level: log.InfoLevel},
		{message: "test3", output: "test3", level: log.DebugLevel},
		{message: "test4", output: "test4", level: log.NoneLevel},
	}

	for _, tt := range testData {
		t.Run(fmt.Sprintf("Write message '%s' correctly", tt.message), func(t *testing.T) {
			output.Reset()
			log.Configure(log.OutputConfig{LogLevel: tt.level})
			bytesCount, err := log.Stdout().Write([]byte(tt.message))
			assert.Equal(t, tt.output, output.String())
			assert.NoError(t, err)
			assert.Equal(t, len([]byte(tt.message)), bytesCount)
		})
	}
}

func mockLogModule() *bytes.Buffer {
	output := bytes.NewBufferString("")
	syslog.SetOutput(output)
	syslog.SetFlags(0)

	return output
}
