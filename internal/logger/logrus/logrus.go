package logrus

import (
	"os"

	"github.com/unknowntpo/todos/internal/logger"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	*logrus.Logger
}

func RegisterLog() error {
	newLogger := &logrusWrapper{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}

	// Let logger.Logger interface use logrus.Logger as implementation.
	logger.SetLogger(newLogger)
	return nil
}

func (lw *logrusWrapper) PrintInfo(message string, properties map[string]string) {
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
	}).Info(message)
}

func (lw *logrusWrapper) PrintError(err error, properties map[string]string) {
	//_, caller, line, _ := runtime.Caller(1)
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
		//"trace":      string(debug.Stack()),
		//"caller":     fmt.Sprintf("%s:%d", caller, line),
		//"trace": fmt.Sprintf("%+v", err),
	}).Error(err)
}

func (lw *logrusWrapper) PrintFatal(err error, properties map[string]string) {
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
	}).Fatal(err)
}
