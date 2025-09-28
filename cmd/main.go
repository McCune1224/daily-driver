package main

import (
	"context"
	handler "daily-driver/internal/handler"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
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

func createDBPool() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func createLogger() (*zap.Logger, error) {
	logFile := filepath.Join(".", "application.log")
	file, err := os.Create(logFile)
	if err != nil {
		return nil, err
	}

	fileSyncer := zapcore.AddSync(file)
	consoleSyncer := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel)
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

	dbPool, err := createDBPool()
	if err != nil {
		logger.Panic("Failed to create database pool: " + err.Error())
	}

	handler := &handler.Handler{
		Logger: logger,
		DBPool: dbPool,
	}

	handler.AttachRoutes(app)
	app.Logger.Fatal(app.Start(":" + getPort()))
}
