package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type LevelHook struct {
	writer    io.Writer
	logLevels []logrus.Level
}

func (hook *LevelHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.writer.Write([]byte(line))
	return err
}

func (hook *LevelHook) Levels() []logrus.Level {
	return hook.logLevels
}

func InitLogger() {
	logFileName := filepath.Join("logs", time.Now().Format("2006-01-02_15-04-05")+"_errors.log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Не удалось открыть файл для логирования: %v", err)
	}

	logrus.SetOutput(io.Discard)

	consoleLogger := logrus.New()
	consoleLogger.SetOutput(os.Stdout)
	consoleLogger.SetLevel(logrus.TraceLevel)

	logrus.AddHook(&LevelHook{
		writer:    logFile,
		logLevels: []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	logrus.AddHook(&LevelHook{
		writer:    os.Stdout,
		logLevels: logrus.AllLevels,
	})

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
