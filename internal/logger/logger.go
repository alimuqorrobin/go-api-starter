// ============================================
// FILE: internal/logger/logger.go
// FIX: Complete logger implementation
// ============================================

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logPath string) *zap.SugaredLogger {
	// Create logs directory if not exists
	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(err)
	}

	// Configure lumberjack for daily rotation with 7 days retention
	w := &lumberjack.Logger{
		Filename:   logPath + "/app.log",
		MaxSize:    100,   // megabytes
		MaxBackups: 7,     // files
		MaxAge:     7,     // days
		Compress:   true,
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// File writer - JSON format untuk logs
	fileWriter := zapcore.AddSync(w)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		zapcore.InfoLevel,
	)

	// Console writer - untuk development readability
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		consoleWriter,
		zapcore.InfoLevel,
	)

	// Combine both cores (file + console)
	multiCore := zapcore.NewTee(fileCore, consoleCore)

	// Create logger with caller info
	logger := zap.New(
		multiCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// Return sugared logger (easier API)
	return logger.Sugar()
}