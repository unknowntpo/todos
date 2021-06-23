package logger

// DefaultLogger
var Log Logger

// Logger represent common interface for logging function
type Logger interface {
	PrintInfo(message string, properties map[string]string)
	PrintError(err error, properties map[string]string)
	PrintFatal(err error, properties map[string]string)
}

func SetLogger(newLogger Logger) {
	Log = newLogger
}
