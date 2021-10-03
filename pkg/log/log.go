package log

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type Level int

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

var (
	ErrorOutput = NewlevelWriter(generalOutput, ErrorLevel)
	DebugOutput = NewlevelWriter(generalOutput, DebugLevel)
	InfoOutput  = NewlevelWriter(generalOutput, InfoLevel)
)

func SetOutput(output io.Writer) {
	generalOutput = output
	ErrorOutput.SetOutput(output)
	DebugOutput.SetOutput(output)
	InfoOutput.SetOutput(output)
}

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

func Debug(params ...interface{}) {
	dbg.Fprintln(DebugOutput, params...)
}

func Debugf(message string, params ...interface{}) {
	dbg.Fprintln(DebugOutput, fmt.Sprintf(message, params...))
}

func Error(params ...interface{}) {
	err.Fprintln(ErrorOutput, params...)
}

func Errorf(message string, params ...interface{}) {
	err.Fprintln(ErrorOutput, fmt.Sprintf(message, params...))
}

func Info(params ...interface{}) {
	inf.Fprintln(InfoOutput, params...)
}

func Infof(message string, params ...interface{}) {
	inf.Fprintln(InfoOutput, fmt.Sprintf(message, params...))
}

func Stdout() io.Writer {
	return generalOutput
}
