package reporter

import (
	"bytes"
	"text/template"

	"github.com/gookit/color"
)

// ConsoleReporter is implementation of Reporter to conminicate with console
type ConsoleReporter struct {
}

// ValidationError prints information about failed rule
func (c *ConsoleReporter) ValidationError(rule, message string) {
	color.Warn.Printf("⛔ [rule: %s] %s", rule, message)
}

// Error prints information about application error
func (c *ConsoleReporter) Error(message string) {
	color.Red.Printf("⛔ %s", message)
}

// Info prints information message
func (c *ConsoleReporter) Info(message string) {
	color.White.Print(message)
}

// Debug prints debug message
func (c *ConsoleReporter) Debug(message string) {
	color.Gray.Printf(message)
}

// PrintGraphics prints template by structure data
func (c *ConsoleReporter) PrintGraphics(content string, data interface{}) {
	tpl, err := template.New(content).Parse(content)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	if err := tpl.Execute(&buff, data); err != nil {
		panic(err)
	}
	color.White.Print(buff.String())
}
