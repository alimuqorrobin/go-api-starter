package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logPath string) *zap.SugaredLogger {
	// Create logs directory
	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(err)
	}

	// Build filename dengan tanggal
	// Format: app-2024-01-16.log
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	logFilename := filepath.Join(logPath, fmt.Sprintf("app-%s.log", dateStr))

	// Configure lumberjack
	w := &lumberjack.Logger{
		Filename:   logFilename,          // app-2024-01-16.log
		MaxSize:    500,                  // Rotate jika >500MB
		MaxBackups: 7,                    // Keep 7 files
		MaxAge:     7,                    // Delete > 7 hari
		Compress:   true,                 // Compress old logs
		LocalTime:  true,                 // Use local time
	}

	// Start background rotation checker
	go startRotationChecker(logPath, w)

	// Encoder config
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

	// File core
	fileWriter := zapcore.AddSync(w)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		zapcore.InfoLevel,
	)

	// Console core
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		consoleWriter,
		zapcore.InfoLevel,
	)

	// Combine cores
	multiCore := zapcore.NewTee(fileCore, consoleCore)

	// Create logger
	logger := zap.New(
		multiCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger.Sugar()
}

// startRotationChecker ensures daily rotation dengan filename berubah
// Ketika hari berubah, tutup file lama dan buat file baru dengan tanggal baru
func startRotationChecker(logPath string, w *lumberjack.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	lastRotateDay := time.Now().Day()

	for range ticker.C {
		now := time.Now()
		currentDay := now.Day()

		// Jika hari berubah, buat file baru dengan tanggal baru
		if currentDay != lastRotateDay {
			// Build filename baru dengan tanggal hari ini
			dateStr := now.Format("2006-01-02")
			newFilename := filepath.Join(logPath, fmt.Sprintf("app-%s.log", dateStr))

			// Update filename di lumberjack
			w.Filename = newFilename

			// Force rotation (tutup file lama, buat file baru)
			err := w.Rotate()
			if err != nil {
				fmt.Printf("[Logger] Rotation error: %v\n", err)
			} else {
				fmt.Printf("[Logger] ✅ Daily rotation at %v - new file: %s\n", 
					now.Format(time.RFC3339), newFilename)
			}

			// Clean old files (lebih dari 7 hari)
			cleanOldLogs(logPath)

			lastRotateDay = currentDay
		}
	}
}

// cleanOldLogs menghapus file log yang lebih dari MaxAge (7 hari)
func cleanOldLogs(logPath string) {
	files, err := os.ReadDir(logPath)
	if err != nil {
		fmt.Printf("[Logger] Error reading log dir: %v\n", err)
		return
	}

	now := time.Now()
	maxAge := 7 * 24 * time.Hour

	for _, file := range files {
		// Skip jika bukan file
		if file.IsDir() {
			continue
		}

		// Skip jika bukan app log file
		filename := file.Name()
		if !isAppLogFile(filename) {
			continue
		}

		// Get file info
		info, err := file.Info()
		if err != nil {
			continue
		}

		// Check age
		age := now.Sub(info.ModTime())
		if age > maxAge {
			// Delete file
			fullPath := filepath.Join(logPath, filename)
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("[Logger] Error deleting %s: %v\n", filename, err)
			} else {
				fmt.Printf("[Logger] ✅ Deleted old log: %s (age: %v days)\n", 
					filename, age.Hours()/24)
			}
		}
	}
}

// isAppLogFile checks jika file adalah app log file
func isAppLogFile(filename string) bool {
	// Match: app-2024-01-16.log atau app-2024-01-16.log.gz
	if len(filename) < 12 {
		return false
	}

	prefix := "app-"
	if filename[:len(prefix)] != prefix {
		return false
	}

	// Check extension
	if filename[len(filename)-4:] == ".log" || filename[len(filename)-7:] == ".log.gz" {
		return true
	}

	return false
}

// ============================================
// CONTOH LOG FILES STRUCTURE
// ============================================

/*

logs/
├── app-2024-01-16.log         # Today (writable)
├── app-2024-01-15.log.gz      # Yesterday (compressed)
├── app-2024-01-14.log.gz      # 2 days old
├── app-2024-01-13.log.gz      # 3 days old
├── app-2024-01-12.log.gz      # 4 days old
├── app-2024-01-11.log.gz      # 5 days old
├── app-2024-01-10.log.gz      # 6 days old
└── app-2024-01-09.log.gz      # 7 days old (will delete next cleanup)

TIMELINE:
- Hari 1 (Jan 16): app-2024-01-16.log (writable)
- Hari 2 (Jan 17 midnight): 
  - app-2024-01-16.log → app-2024-01-16.log.gz (compressed)
  - Create app-2024-01-17.log (baru)
- Hari 8 (Jan 23): 
  - Cleanup: app-2024-01-16.log.gz (older than 7 days) → DELETE
  - Keep: app-2024-01-17 sampai app-2024-01-23

*/