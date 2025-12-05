package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sug *zap.SugaredLogger

// Init configures a global sugared logger based on a level string (debug/info/warn/error).
// If filePath is not empty, logs will be written to both stdout and the specified file.
// It returns an error if the underlying logger cannot be built.
func Init(level string, filePath string) error {
	lvl := zapcore.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn", "warning":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	default:
		lvl = zapcore.InfoLevel
	}

	// Configure output paths
	outputPaths := []string{"stdout"}
	if filePath != "" {
		outputPaths = append(outputPaths, filePath)
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	// replace global sug
	sug = zl.Sugar()
	return nil
}

// Sugar returns the global *zap.SugaredLogger. It will lazily initialize an
// info-level logger if Init wasn't called.
func Sugar() *zap.SugaredLogger {
	if sug == nil {
		_ = Init("info", "")
	}
	return sug
}

// Logger returns the underlying *zap.Logger. It will lazily initialize if needed.
func Logger() *zap.Logger {
	if sug == nil {
		_ = Init("info", "")
	}
	return sug.Desugar()
}

// Sync flushes any buffered log entries.
func Sync() error {
	if sug == nil {
		return nil
	}
	return Logger().Sync()
}

// Debug logs a debug message with optional key-value pairs.
func Debug(msg string, keysAndValues ...interface{}) {
	Sugar().Debugw(msg, keysAndValues...)
}

// Info logs an info message with optional key-value pairs.
func Info(msg string, keysAndValues ...interface{}) {
	Sugar().Infow(msg, keysAndValues...)
}

// Warn logs a warning message with optional key-value pairs.
func Warn(msg string, keysAndValues ...interface{}) {
	Sugar().Warnw(msg, keysAndValues...)
}

// Error logs an error message with optional key-value pairs.
func Error(msg string, keysAndValues ...interface{}) {
	Sugar().Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message with optional key-value pairs and exits.
func Fatal(msg string, keysAndValues ...interface{}) {
	Sugar().Fatalw(msg, keysAndValues...)
}
