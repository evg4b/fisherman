package log

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	NoneLevel
)

var generalOutput io.Writer = os.Stdout

var (
	err = color.New(color.FgRed).Add(color.Bold)
	dbg = color.New(color.FgYellow)
	inf = color.New(color.FgWhite)
)

var (
	ErrorOutput = NewLevelWriter(generalOutput, ErrorLevel)
	DebugOutput = NewLevelWriter(generalOutput, DebugLevel)
	InfoOutput  = NewLevelWriter(generalOutput, InfoLevel)
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

func Debug(params ...any) {
	dbg.Fprintln(DebugOutput, params...)
}

func Debugf(message string, params ...any) {
	dbg.Fprintln(DebugOutput, fmt.Sprintf(message, params...))
}

func Error(params ...any) {
	err.Fprintln(ErrorOutput, params...)
}

func Errorf(message string, params ...any) {
	err.Fprintln(ErrorOutput, fmt.Sprintf(message, params...))
}

func Info(params ...any) {
	inf.Fprintln(InfoOutput, params...)
}

func Infof(message string, params ...any) {
	inf.Fprintln(InfoOutput, fmt.Sprintf(message, params...))
}

func Stdout() io.Writer {
	return generalOutput
}
