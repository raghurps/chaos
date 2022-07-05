package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	cfg := zap.NewDevelopmentConfig()

	// Set log level from env variable
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		cfg.Level.SetLevel(zapcore.DebugLevel)
	case "warn":
		cfg.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	default:
		cfg.Level.SetLevel(zapcore.InfoLevel)
	}

	// Log everything to STDOUT
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	cfg.Encoding = "json"

	Logger, _ = cfg.Build()
}
