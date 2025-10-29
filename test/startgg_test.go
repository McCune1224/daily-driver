package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"daily-driver/internal/handler"
	"daily-driver/internal/routes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func TestStartGGPanelRoute(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()

	// Create a nil DB pool for this test (we're not using DB in the StartGG handler)
	var dbPool *pgxpool.Pool

	// Create handler
	h := handler.NewHandler(logger, dbPool)

	// Create Echo instance
	e := echo.New()

	// Attach routes
	h.AttachRoutes(e)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, routes.StartGGBase, nil)
	rec := httptest.NewRecorder()

	// Create a new context
	c := e.NewContext(req, rec)

	// Set the path to match the route
	c.SetPath(routes.StartGGBase)

	// Call the handler directly
	err := h.RenderPanelStartGG(c)

	// Assert no error
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Assert status code
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got: %d", rec.Code)
	}

	// Assert response contains expected content
	body := rec.Body.String()
	if len(body) == 0 {
		t.Fatal("Expected non-empty response body")
	}

	// Check for Start.GG content - could be either success or error case
	expectedStrings := []string{
		"START.GG_ESPORTS.SYS",
	}

	// Check if it's the success case (with emoji) or error case
	successIndicators := []string{"🎮", "RECENT_TOURNAMENTS"}
	errorIndicators := []string{"ERROR", "API"}

	foundSuccess := false
	foundError := false

	for _, indicator := range successIndicators {
		if contains(body, indicator) {
			foundSuccess = true
			break
		}
	}

	for _, indicator := range errorIndicators {
		if contains(body, indicator) {
			foundError = true
			break
		}
	}

	// Should have either success content or error content
	if !foundSuccess && !foundError {
		t.Fatalf("Expected response to contain either success indicators (%v) or error indicators (%v), but it contained neither", successIndicators, errorIndicators)
	}

	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Fatalf("Expected response to contain '%s', but it didn't", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
