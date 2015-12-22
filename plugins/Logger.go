package plugins

type Logger interface {
	Error(msg string, params ...interface{})
	Warn(msg string, params ...interface{})
}
