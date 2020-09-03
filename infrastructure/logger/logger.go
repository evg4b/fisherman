package logger

type Logger interface {
	Debug(message string)
	Error(message string)
	Info(message string)
}
