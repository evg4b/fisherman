package logger

import "io"

// LogLevel is type for determinate log level for console logger
type LogLevel int

// Available log levels
const (
	Debug LogLevel = iota
	Info
	Error
	None
)

// Logger is interface for information logger
type Logger interface {
	io.Writer
	// Debug prints diagnostic message to output.
	// Output can be skepped when log level is `Info`, `Error` or `None`.
	Debug(params ...interface{})
	// Debugf prints diagnostic message with formatting to output.
	// Output can be skepped when log level is `Info`, `Error` or `None`.
	Debugf(message string, params ...interface{})
	// Error prints error message to output.
	// Output can be skepped when log level is `None`.
	Error(params ...interface{})
	// Errorf prints error message with formatting to output.
	// Output can be skepped when log level is `None`.
	Errorf(message string, params ...interface{})
	// Info prints information message to output
	// Output can be skepped when log level is `Error` or `None`.
	Info(params ...interface{})
	// Infof prints information message with formatting to output
	// Output can be skepped when log level is `Error` or `None`.
	Infof(message string, params ...interface{})
}
