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

func (logger *ConsoleLogger) Debug(params ...interface{}) {
	if logger.level >= Debug {
		logger.debugPrinter.Println(params)
	}
}

func (logger *ConsoleLogger) Debugf(message string, params ...interface{}) {
	if logger.level >= Debug {
		logger.debugPrinter.Printf(message, params...)
	}
}

func (logger *ConsoleLogger) Error(params ...interface{}) {
	if logger.level >= Error {
		logger.errorPrinter.Println(params)
	}
}

func (logger *ConsoleLogger) Errorf(message string, params ...interface{}) {
	if logger.level >= Error {
		logger.errorPrinter.Printf(message, params...)
	}
}

func (logger *ConsoleLogger) Info(params ...interface{}) {
	if logger.level >= Info {
		logger.infoPrinter.Println(params)
	}
}

func (logger *ConsoleLogger) Infof(message string, params ...interface{}) {
	if logger.level >= Info {
		logger.infoPrinter.Printf(message, params...)
	}
}
