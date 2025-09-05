package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/", "static")

	hub := NewHub()
	go hub.Run()

	// WebSocket endpoint
	e.GET("/ws", func(c echo.Context) error {
		return HandleWebSocket(hub, c)
	})

	// API routes
	api := e.Group("/api")
	api.GET("/weather", GetWeather)
	api.GET("/garmin", GetGarminData)
	api.GET("/system", GetSystemInfo)

	log.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
