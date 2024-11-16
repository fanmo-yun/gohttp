package logger

import (
	"fmt"
	"gohttpd/utils"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(lc utils.LoggerConfig) {
	switch lc.Out {
	case "console":
		NewStdOutLogger(lc.Level)
	case "file":
		NewFileLogger(lc.Level)
	default:
		fmt.Fprintf(os.Stdout, "gohttp: Log Output Error type\n")
		os.Exit(1002)
	}
}

func NewStdOutLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		fmt.Fprintf(os.Stdout, "gohttp: Log Level Init Fatal: %v\n", err)
		os.Exit(1002)
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stdout, "gohttp: Log Build Fatal: %v\n", err)
		os.Exit(1002)
	}

	zap.ReplaceGlobals(logger)
}

func NewFileLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		fmt.Fprintf(os.Stdout, "gohttp: Log Level Init Fatal: %v\n", err)
		os.Exit(1002)
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	logFile := fmt.Sprintf("%s.log", timestamp)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join("log", logFile),
		MaxSize:    128,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		logLevel,
	)

	logger := zap.New(core)

	zap.ReplaceGlobals(logger)
}
