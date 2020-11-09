package log

import (
	"io"
	"io/ioutil"
	"log"

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

var level = InfoLevel

// Configure configures logger
func Configure(config OutputConfig) {
	if !config.Colors {
		err.DisableColor()
		dbg.DisableColor()
		inf.DisableColor()
	}
	level = config.LogLevel
}

// Debug prints diagnostic message to output with debug styles (Yealow font) to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Debug(params ...interface{}) {
	if level <= DebugLevel {
		log.Println(dbg.Sprint(params...))
	}
}

// Debugf prints diagnostic message to output with debug styles (Yealow font) and formatting to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Debugf(message string, params ...interface{}) {
	if level <= DebugLevel {
		log.Println(dbg.Sprintf(message, params...))
	}
}

// Error prints error to console with error styles (Bold red font) to output.
// Output can be skepped when log level parameter is None
// Color styles can be omitted when color paramerter is false.
func Error(params ...interface{}) {
	if level <= ErrorLevel {
		log.Println(err.Sprint(params...))
	}
}

// Errorf prints error message with error styles (Bold red font) and formatting to output
// Output can be skepped when log level is `None`.
// Color styles can be omitted when color paramerter is false.s
func Errorf(message string, params ...interface{}) {
	if level <= ErrorLevel {
		log.Println(err.Sprintf(message, params...))
	}
}

// Info prints information message with information styles (withe font) to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Info(params ...interface{}) {
	if level <= InfoLevel {
		log.Println(inf.Sprint(params...))
	}
}

// Infof prints information message with information styles (withe font) and formatting to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func Infof(message string, params ...interface{}) {
	if level <= InfoLevel {
		log.Println(inf.Sprintf(message, params...))
	}
}

// Writer return output io.Writer object
func Writer() io.Writer {
	if level <= InfoLevel {
		return log.Writer()
	}

	return ioutil.Discard
}

// Stdout return output io.Writer object withoud level handling
func Stdout() io.Writer {
	return log.Writer()
}
