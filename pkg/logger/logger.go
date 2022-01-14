package logger

import (
	log "github.com/sirupsen/logrus"
)

var logger log.Logger

func Init(level log.Level) {
	logger = *log.New()
	logger.SetLevel(level)
}

func Success(message ...string) {
	if len(message) == 0 {
		logger.Debug("[SUCCESS]")
		return
	}
	logger.Debug(message[0] + "[SUCCESS]")
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Error(message string, err error) {
	logger.Error(message+" error = ", err)
}
