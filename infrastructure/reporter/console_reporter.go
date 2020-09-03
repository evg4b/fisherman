package reporter

import (
	"bytes"
	"fmt"
	"text/template"
)

// ConsoleReporter is implementation of Reporter to conminicate with console
type ConsoleReporter struct {
}

// NewConsoleReporter returns new instance of ConsoleReporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// ValidationError prints information about failed rule
func (c *ConsoleReporter) ValidationError(rule, message string) {
	fmt.Printf("⛔ [rule: %s] %s", rule, message)
}

// Error prints information about application error
func (c *ConsoleReporter) Error(message string) {
	fmt.Printf("⛔ %s", message)
}

// Info prints information message
func (c *ConsoleReporter) Info(message string) {
	fmt.Print(message)
}

// Debug prints debug message
func (c *ConsoleReporter) Debug(message string) {
	fmt.Printf(message)
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
	fmt.Print(buff.String())
}
