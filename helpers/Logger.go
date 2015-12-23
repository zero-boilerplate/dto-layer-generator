package helpers

type Logger interface {
	Error(msg string, params ...interface{})
	Warn(msg string, params ...interface{})
	Info(msg string, params ...interface{})
	Debug(msg string, params ...interface{})
}
