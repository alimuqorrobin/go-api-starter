package logger

import (
    "os"

    "golang-api-starter/config"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
    *zap.SugaredLogger
}

func NewLogger(cfg *config.Config) *Logger {
    // Ensure logs directory exists
    if err := os.MkdirAll("logs", 0755); err != nil {
        panic(err)
    }

    // Configure lumberjack for daily log rotation
    writer := &lumberjack.Logger{
        Filename:   cfg.LogFilePath,
        MaxSize:    cfg.LogMaxSize,    // megabytes
        MaxBackups: cfg.LogMaxBackups, // number of backups
        MaxAge:     cfg.LogMaxAge,     // days
        Compress:   true,              // compress old log files
    }

    // Configure encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "timestamp"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.StacktraceKey = ""
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

    // Create core with multiple outputs
    core := zapcore.NewTee(
        // File output
        zapcore.NewCore(
            zapcore.NewJSONEncoder(encoderConfig),
            zapcore.AddSync(writer),
            getLogLevel(cfg.LogLevel),
        ),
        // Console output
        zapcore.NewCore(
            zapcore.NewConsoleEncoder(encoderConfig),
            zapcore.AddSync(os.Stdout),
            getLogLevel(cfg.LogLevel),
        ),
    )

    // Build logger
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    return &Logger{
        SugaredLogger: logger.Sugar(),
    }
}

func getLogLevel(level string) zapcore.Level {
    switch level {
    case "debug":
        return zapcore.DebugLevel
    case "info":
        return zapcore.InfoLevel
    case "warn":
        return zapcore.WarnLevel
    case "error":
        return zapcore.ErrorLevel
    default:
        return zapcore.InfoLevel
    }
}
```

**Penjelasan:**
- **Zap logger** = High-performance structured logging
- **Lumberjack** = Log rotation otomatis per hari
- **Multiple outputs** = Console + File
- JSON format di file, console format di terminal

**Log Rotation:**
```
// logs/
// ├── app.log              ← Current log
// ├── app-2024-01-01.log   ← Rotated (compressed)
// ├── app-2024-01-02.log
// └── app-2024-01-03.log

// Config:
// - MaxSize: 100MB per file
// - MaxBackups: 30 files
// - MaxAge: 90 days
// - Auto compress old logs
// ```

// **Log Levels:**
// ```
// DEBUG → Detailed info untuk debugging
// INFO  → General info (requests, etc)
// WARN  → Warning messages
// ERROR → Error messages + stack trace

// log := logger.NewLogger(cfg)

// log.Info("Server started")
// log.Infow("Request processed", "path", "/users", "status", 200)
// log.Error("Database error", "error", err)
// log.Fatal("Critical error", "error", err) // Exit app