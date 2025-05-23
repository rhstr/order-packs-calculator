package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Execute runs the application.
func Execute() {
	initializeLogger()
	defer zap.L().Sync()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	err := e.Start(":8080")
	if err != nil {
		zap.L().Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

// initializeLogger sets up the global logger for the application.
func initializeLogger() {
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
