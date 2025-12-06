package logger

import (
    "io"
    "time"

    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger
type Logger struct {
    *zap.Logger
}

// NewLogger creates a zap logger writing to daily rotated files under logDir.
func NewLogger(logDir string, rotationDays int, maxAgeDays int, level string) (*Logger, error) {
    writer, err := rotatelogs.New(
        logDir + "/app.%Y-%m-%d.log",
        rotatelogs.WithRotationTime(time.Duration(rotationDays)*24*time.Hour),
        rotatelogs.WithMaxAge(time.Duration(maxAgeDays)*24*time.Hour),
    )
    if err != nil {
        return nil, err
    }
    return newZapWithWriter(writer, level), nil
}

func newZapWithWriter(w io.Writer, level string) (*Logger, error) {
    encCfg := zap.NewProductionEncoderConfig()
    encCfg.TimeKey = "ts"
    encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

    atom := zap.NewAtomicLevel()
    switch level {
    case "debug":
        atom.SetLevel(zap.DebugLevel)
    case "info":
        atom.SetLevel(zap.InfoLevel)
    case "warn":
        atom.SetLevel(zap.WarnLevel)
    default:
        atom.SetLevel(zap.InfoLevel)
    }

    core := zapcore.NewCore(zapcore.NewJSONEncoder(encCfg),
        zapcore.AddSync(w),
        atom,
    )

    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    return &Logger{logger}, nil
}

// F creates a field for logging
func F(k string, v interface{}) zap.Field { return zap.Any(k, v) }

// GoSafe runs a goroutine with recover + logging, to emulate try/catch in goroutines.
func (l *Logger) GoSafe(f func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                l.Error("panic in goroutine", F("panic", r))
            }
        }()
        f()
    }()
}
