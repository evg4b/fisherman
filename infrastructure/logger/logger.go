package logger

type Logger interface {
	Debug(params ...interface{})
	Debugf(message string, params ...interface{})
	Error(params ...interface{})
	Errorf(message string, params ...interface{})
	Info(params ...interface{})
	Infof(message string, params ...interface{})
}
