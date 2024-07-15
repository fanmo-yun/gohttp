package logger

import (
	"fmt"
	"gohttp/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(lc utils.LoggerConfig) {
	switch lc.Out {
	case "console":
		NewStdOutLogger(lc.Level)
	case "file":
		NewFileLogger(lc.Level)
	default:
		fmt.Fprintf(os.Stdout, "Log Out Error type\n")
		os.Exit(1)
	}
}

func NewStdOutLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		fmt.Fprintf(os.Stdout, "Log Level Init Fatal: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stdout, "Log Build Fatal: %v\n", err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)
}

func NewFileLogger(l string) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(l))
	if err != nil {
		fmt.Fprintf(os.Stdout, "Log Level Init Fatal: %v\n", err)
		os.Exit(1)
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"logfile.log"},
		ErrorOutputPaths: []string{"logfile.err"},
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Log Build Fatal: %v\n", err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)
}
