package logger

import (
    "fmt"
    "io"
    "os"
    "time"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
)

type Logger struct {
    *logrus.Logger
}

var std *Logger

func NewLogger(logDir string, rotationDays int, maxAgeDays int, level string) (*Logger, error) {
    if err := os.MkdirAll(logDir, 0o755); err != nil {
        return nil, err
    }
    path := fmt.Sprintf("%s/app", logDir)
    writer, err := rotatelogs.New(path+"-%Y-%m-%d.log", rotatelogs.WithRotationTime(time.Duration(rotationDays)*24*time.Hour), rotatelogs.WithMaxAge(time.Duration(maxAgeDays)*24*time.Hour))
    if err != nil {
        return nil, err
    }
    lg := logrus.New()
    lg.SetFormatter(&logrus.JSONFormatter{TimestampFormat:time.RFC3339})
    switch level {
    case "debug":
        lg.SetLevel(logrus.DebugLevel)
    case "warn":
        lg.SetLevel(logrus.WarnLevel)
    default:
        lg.SetLevel(logrus.InfoLevel)
    }
    mw := io.MultiWriter(os.Stdout, writer)
    lg.SetOutput(mw)
    hook := lfshook.NewHook(lfshook.WriterMap{
        logrus.InfoLevel: writer,
        logrus.ErrorLevel: writer,
        logrus.WarnLevel: writer,
        logrus.DebugLevel: writer,
    }, &logrus.JSONFormatter{})
    lg.AddHook(hook)
    std = &Logger{lg}
    return std, nil
}

func F(k string, v interface{}) map[string]interface{} {
    return map[string]interface{}{k: v}
}
