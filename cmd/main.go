package main

import (
	handler "daily-driver/internal/handler"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DevPort = "8080"

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return DevPort
	}
	return port
}

func createLogger() (*zap.Logger, error) {
	logFile := filepath.Join(".", "application.log")
	file, err := os.Create(logFile)
	if err != nil {
		return nil, err
	}

	// File logging syncer
	fileSyncer := zapcore.AddSync(file)

	// Console logging syncer
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// Encoder configuration
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Core for file and console logging
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel)

	// Combine cores
	combinedCore := zapcore.NewTee(fileCore, consoleCore)

	return zap.New(combinedCore, zap.AddStacktrace(zap.FatalLevel)), nil
}

func main() {
	// Initialize Zap logger to log to both file and standard output
	logger, err := createLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Create Echo instance
	app := echo.New()

	// Replace Echo's default logger with Zap
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("Incoming request",
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
			)
			return next(c)
		}
	})

	app.Use(middleware.Recover())

	app.Static("/static", "web/static")

	handler := &handler.Handler{
		Logger: logger,
	}
	handler.AttachRoutes(app)
	app.Logger.Fatal(app.Start(":" + getPort()))
}
