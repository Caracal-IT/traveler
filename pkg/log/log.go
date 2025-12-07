package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"traveler/pkg/config"
)

var sug *zap.SugaredLogger

// Init configures a global sugared logger based on a level string (debug/info/warn/error).
// If filePath is not empty, logs will be written to both stdout and the specified file.
// If esCfg is provided and Enabled, logs will also be shipped to Elasticsearch.
// It returns an error if the underlying logger cannot be built.
func Init(level string, filePath string, esCfg *config.ElasticLogConfig) error {
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

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
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
	}

	// Build base writer(s)
	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if filePath != "" {
		if f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			writers = append(writers, zapcore.AddSync(f))
		}
	}
	baseCore := zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.NewMultiWriteSyncer(writers...), lvl)

	core := baseCore

	// Optionally add Elasticsearch core
	if esCfg != nil && esCfg.Enabled && esCfg.URL != "" && esCfg.Index != "" {
		esSyncer := newElasticsearchSyncer(esCfg.URL, esCfg.Index)
		esCore := zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), esSyncer, lvl)
		core = zapcore.NewTee(core, esCore)
	}

	zl := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	sug = zl.Sugar()
	return nil
}

// Sugar returns the global *zap.SugaredLogger. It will lazily initialize an
// info-level logger if Init wasn't called.
func Sugar() *zap.SugaredLogger {
	if sug == nil {
		_ = Init("info", "", nil)
	}
	return sug
}

// Logger returns the underlying *zap.Logger. It will lazily initialize if needed.
func Logger() *zap.Logger {
	if sug == nil {
		_ = Init("info", "", nil)
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
