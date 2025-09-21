package main

import (
	handler "daily-driver/internal/handler"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const DevPort = "8080"

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return DevPort
	}
	return port
}

func main() {
	// Initialize Zap logger with development (colored) configuration
	logger, err := zap.NewDevelopment(
		zap.AddStacktrace(zap.FatalLevel),
		zap.Development(),
	) // Enables colored console-friendly logging
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
