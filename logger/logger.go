package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logFile := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	Log.Level = logrus.ErrorLevel
}

func Info(message string, args ...interface{}) {
	Log.Infof(message, args...)
}

func Error(message string, args ...interface{}) {
	Log.Errorf(message, args...)
}

func Fatal(message string, args ...interface{}) {
	Log.Fatalf(message, args...)
}

func Warn(message string, args ...interface{}) {
	Log.Warnf(message, args...)
}
