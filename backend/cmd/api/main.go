package main

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	// Initialize Zap logger to human-readable development setup
	var err error
	logger, err = zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()

	logger.Info("Logger initialized")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found")
	}

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	// Initialize Echo
	e := echo.New()

	// Middleware
	// Use Zap logger as middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("Request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
			)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Serve static files from SvelteKit build (production)
	e.Static("/", "./frontend/build")

	// API routes
	api := e.Group("/api")

	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"message": "Aperture Science Data Portal is operational",
		})
	})

	// Activity routes
	api.GET("/activities", getActivities)
	api.POST("/activities", createActivity)

	// Tournament routes
	api.GET("/tournaments", getTournaments)
	api.POST("/tournaments", createTournament)

	// Art routes
	api.GET("/art/random", getRandomArtwork)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", zap.String("port", port))
	e.Logger.Fatal(e.Start(":" + port))
}

// Placeholder handlers - implement these in handlers package
func getActivities(c echo.Context) error {
	return c.JSON(200, []map[string]interface{}{
		{
			"id":               1,
			"activity_type":    "run",
			"distance_meters":  5000.0,
			"duration_seconds": 1800,
			"activity_date":    "2025-11-01T08:00:00Z",
		},
	})
}

func createActivity(c echo.Context) error {
	var activity map[string]interface{}
	if err := c.Bind(&activity); err != nil {
		logger.Error("Failed to bind request payload", zap.Error(err))
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	return c.JSON(201, activity)
}

func getTournaments(c echo.Context) error {
	return c.JSON(200, []map[string]interface{}{
		{
			"id":              1,
			"tournament_name": "Genesis 10",
			"game":            "Super Smash Bros. Ultimate",
			"placement":       25,
			"tournament_date": "2025-10-15",
		},
	})
}

func createTournament(c echo.Context) error {
	var tournament map[string]interface{}
	if err := c.Bind(&tournament); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	return c.JSON(201, tournament)
}

func getRandomArtwork(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"title":        "A Sunday on La Grande Jatte",
		"artist":       "Georges Seurat",
		"date_display": "1884-1886",
		"image_url":    "https://www.artic.edu/iiif/2/...",
	})
}
