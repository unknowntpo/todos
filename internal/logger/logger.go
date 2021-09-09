package logger

// Logger represent common interface for logging function
type Logger interface {
	PrintInfo(message string, properties map[string]interface{})
	PrintError(err error, properties map[string]interface{})
	PrintFatal(err error, properties map[string]interface{})
}
