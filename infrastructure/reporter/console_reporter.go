package reporter

import "github.com/gookit/color"

const (
	Debug = 0
)

type ConsoleReporter struct {
}

func (c *ConsoleReporter) ValidationError(rule, message string) {
	color.Warn.Printf("⛔ %s", message)
}

func (c *ConsoleReporter) Error(message string) {
	color.Red.Printf("⛔ %s", message)
}

func (c *ConsoleReporter) Info(message string) {
	color.White.Printf(message)
	color.Red.Printf(message)
	color.Error.Printf(message)
	color.Warn.Printf("⛔ %s", message)
}

func (c *ConsoleReporter) Debug(message string) {
	color.Gray.Printf(message)
}
