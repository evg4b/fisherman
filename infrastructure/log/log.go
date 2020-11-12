package log

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Level is type for determinate log level for console logger
type Level int

// Available log levels
const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	NoneLevel
)

var (
	err *color.Color = color.New(color.FgRed).Add(color.Bold)
	dbg *color.Color = color.New(color.FgYellow)
	inf *color.Color = color.New(color.FgWhite)
)

var generalOutput io.Writer = os.Stdout
var ErrorOutput = NewlevelWriter(generalOutput, ErrorLevel)
var DebugOutput = NewlevelWriter(generalOutput, DebugLevel)
var InfoOutput = NewlevelWriter(generalOutput, InfoLevel)

func SetOutput(output io.Writer) {
	generalOutput = output
	ErrorOutput.SetOutput(output)
	DebugOutput.SetOutput(output)
	InfoOutput.SetOutput(output)
}

// Configure configures logger
func Configure(config OutputConfig) {
	if !config.Colors {
		err.DisableColor()
		dbg.DisableColor()
		inf.DisableColor()
	}
	ErrorOutput.SetLevel(config.LogLevel)
	DebugOutput.SetLevel(config.LogLevel)
	InfoOutput.SetLevel(config.LogLevel)
}

// Debug prints diagnostic message to output with debug styles (Yealow font) to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Debug(params ...interface{}) {
	dbg.Fprintln(DebugOutput, params...)
}

// Debugf prints diagnostic message to output with debug styles (Yealow font) and formatting to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Debugf(message string, params ...interface{}) {
	dbg.Fprintln(DebugOutput, fmt.Sprintf(message, params...))
}

// Error prints error to console with error styles (Bold red font) to output.
// Output can be skepped when log level parameter is None
// Color styles can be omitted when color paramerter is false.
func Error(params ...interface{}) {
	err.Fprintln(ErrorOutput, params...)
}

// Errorf prints error message with error styles (Bold red font) and formatting to output
// Output can be skepped when log level is `None`.
// Color styles can be omitted when color paramerter is false.s
func Errorf(message string, params ...interface{}) {
	err.Fprintln(ErrorOutput, fmt.Sprintf(message, params...))
}

// Info prints information message with information styles (withe font) to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Info(params ...interface{}) {
	inf.Fprintln(InfoOutput, params...)
}

// Infof prints information message with information styles (withe font) and formatting to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Infof(message string, params ...interface{}) {
	inf.Fprintln(InfoOutput, fmt.Sprintf(message, params...))
}

// Stdout return output io.Writer object withoud level handling
func Stdout() io.Writer {
	return generalOutput
}
