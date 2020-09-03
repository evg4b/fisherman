package logger

import (
	"github.com/fatih/color"
)

type LogLevel int

// Available log levels
const (
	None LogLevel = iota
	Debug
	Info
	Error
)

type Printer = func(format string, a ...interface{})

type ConsoleLogger struct {
	level        LogLevel
	errorPrinter color.Color
	debugPrinter color.Color
	infoPrinter  color.Color
}

func NewConsoleLooger(configuration OutputConfig) *ConsoleLogger {
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

func (logger *ConsoleLogger) Debug(message string) {
	if logger.level >= Debug {
		logger.debugPrinter.Println(message)
	}
}

func (logger *ConsoleLogger) Error(message string) {
	if logger.level >= Error {
		logger.errorPrinter.Println(message)
	}
}

func (logger *ConsoleLogger) Info(message string) {
	if logger.level >= Info {
		logger.infoPrinter.Println(message)
	}
}
