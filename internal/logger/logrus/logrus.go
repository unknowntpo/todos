package logrus

import (
	"os"

	"github.com/unknowntpo/todos/internal/logger"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	*logrus.Logger
}

func RegisterLog() logger.Logger {
	return &logrusWrapper{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}
}

func (lw *logrusWrapper) PrintInfo(message string, properties map[string]string) {
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
	}).Info(message)
}

func (lw *logrusWrapper) PrintError(err error, properties map[string]string) {
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
	}).Error(err)
}

func (lw *logrusWrapper) PrintFatal(err error, properties map[string]string) {
	lw.Logger.WithFields(logrus.Fields{
		"properties": properties,
	}).Fatal(err)
}
