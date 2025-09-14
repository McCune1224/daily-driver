package main

import (
	handler "daily-driver/internal/handler"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	app := echo.New()

	app.Use(middleware.Recover())
	app.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, error=[${error}]\n",
		},
	))

	//static files:
	app.Static("/static", "web/static")

	handler := &handler.Handler{}
	handler.AttachRoutes(app)
	// start server:
	app.Logger.Fatal(app.Start(":" + getPort()))
}
