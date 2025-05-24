package cmd

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Execute runs the application.
func Execute() {
	initializeZapLogger()
	defer zap.L().Sync()

	cfg := serviceContainer.getConfig()
	server := serviceContainer.getHTTPServer()

	err := server.Start(":" + cfg.Port)
	if err != nil {
		zap.L().Fatal("Failed to run HTTP server", zap.Error(err))
	}
}

// initializeZapLogger sets up the global logger for the application.
func initializeZapLogger() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeCaller = func(_ zapcore.EntryCaller, _ zapcore.PrimitiveArrayEncoder) {} // omit caller

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.EncoderConfig = encoderCfg
	loggerCfg.Encoding = "console"

	logger, err := loggerCfg.Build()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize logger: %s\n", err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)
}
