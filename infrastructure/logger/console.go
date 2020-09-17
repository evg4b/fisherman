package logger

import (
	"fmt"

	"github.com/fatih/color"
)

// ConsoleLogger is base structure for storage data for logger.
// This structure implements `logger.Logger` and `io.Writer` interfaces.
type ConsoleLogger struct {
	level        LogLevel
	errorPrinter color.Color
	debugPrinter color.Color
	infoPrinter  color.Color
}

// NewConsoleLogger creates new instance of ConsoleLogger by passed configuration
func NewConsoleLogger(configuration OutputConfig) *ConsoleLogger {
	logger := ConsoleLogger{
		level:        configuration.LogLevel,
		errorPrinter: *color.New(color.FgRed).Add(color.Bold),
		debugPrinter: *color.New(color.FgYellow),
		infoPrinter:  *color.New(color.FgWhite),
	}

	if !configuration.Colors {
		logger.errorPrinter.DisableColor()
		logger.debugPrinter.DisableColor()
		logger.infoPrinter.DisableColor()
	}

	return &logger
}

// Debug prints diagnostic message to output with debug styles (Yealow font) to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func (logger *ConsoleLogger) Debug(params ...interface{}) {
	if logger.level <= Debug {
		logger.debugPrinter.Println(params...)
	}
}

// Debugf prints diagnostic message to output with debug styles (Yealow font) and formatting to output.
// Output can be skepped when log level is `Info`, `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func (logger *ConsoleLogger) Debugf(message string, params ...interface{}) {
	if logger.level <= Debug {
		logger.debugPrinter.Println(fmt.Sprintf(message, params...))
	}
}

// Error prints error to console with error styles (Bold red font) to output.
// Output can be skepped when log level parameter is None
// Color styles can be omitted when color paramerter is false.
func (logger *ConsoleLogger) Error(params ...interface{}) {
	if logger.level <= Error {
		logger.errorPrinter.Println(params...)
	}
}

// Errorf prints error message with error styles (Bold red font) and formatting to output
// Output can be skepped when log level is `None`.
// Color styles can be omitted when color paramerter is false.s
func (logger *ConsoleLogger) Errorf(message string, params ...interface{}) {
	if logger.level <= Error {
		logger.errorPrinter.Println(fmt.Sprintf(message, params...))
	}
}

// Info prints information message with information styles (withe font) to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func (logger *ConsoleLogger) Info(params ...interface{}) {
	if logger.level <= Info {
		logger.infoPrinter.Println(params...)
	}
}

// Infof prints information message with information styles (withe font) and formatting to output.
// Output can be skepped when log level is `Error` or `None`.
// Color styles can be omitted when color paramerter is false.
func (logger *ConsoleLogger) Infof(message string, params ...interface{}) {
	if logger.level <= Info {
		logger.infoPrinter.Println(fmt.Sprintf(message, params...))
	}
}

// Write is implementation io.Writer interface to comunicate with information output.
// Output from this method can not be skipped by log level.
func (logger *ConsoleLogger) Write(p []byte) (n int, err error) {
	return logger.infoPrinter.Print(string(p))
}
