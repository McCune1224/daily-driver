package handler

import (
	"daily-driver/internal/routes"
	"daily-driver/internal/service"
	"daily-driver/web/static/templates"
	"daily-driver/web/static/templates/panel"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

// Helper functions for mock data
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// TODO: Add DB stuff...
type Handler struct {
	PanelHandlers []func(echo.Context) error
	Logger        *zap.Logger
	DBPool        *pgxpool.Pool
	// GarminData []*proto.FIT
}

func NewHandler(logger *zap.Logger, dbPool *pgxpool.Pool) *Handler {
	h := &Handler{
		Logger: logger,
		DBPool: dbPool,
	}
	h.PanelHandlers = []func(echo.Context) error{
		// func(c echo.Context) error { return h.RenderPanelWeather(c) },
		func(c echo.Context) error { return h.RenderGarminPanel(c) },
		func(c echo.Context) error { return h.RenderPanelStartGG(c) },
		// func(c echo.Context) error { return h.RenderGarminPanel(c) },
	}
	h.Logger.Info("Handler initialized")
	h.Logger.Info("Handler initialized", zap.Int("num_panel_handlers", len(h.PanelHandlers)))
	return h
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func (h *Handler) RenderIndex(c echo.Context) error {
	return Render(c, 200, templates.Index())
}

func (h *Handler) AttachRoutes(e *echo.Echo) {
	e.GET(routes.Root, h.RenderIndex)
	e.GET(routes.WeatherAPI, h.RenderPanelWeather)

	//garmin stuff
	e.GET(routes.GarminBase, h.RenderGarminUploadPage)
	e.POST(routes.GarminUpload, h.UploadGarminFile)

	// Panel routes
	e.GET(routes.PanelBase, h.RenderPanels)
	e.GET(routes.PanelIndex, h.UpdatePanelIndex)

	// Art routes
	e.GET(routes.ArtRandomAPI, h.RenderPanelArtwork)

	// StartGG routes
	e.GET(routes.StartGGBase, h.RenderPanelStartGG)
}

// RenderPanelStartGG renders the Start.GG esports panel
func (h *Handler) RenderPanelStartGG(c echo.Context) error {
	// Create Start.GG client
	client, err := service.NewStartGGClient()
	if err != nil {
		h.Logger.Error("Failed to create Start.GG client", zap.Error(err))
		return Render(c, 200, panel.PanelStartGG(&panel.StartGGPanelData{
			UserSlug:   "mckusa", // Default user - could be made configurable
			LastUpdate: time.Now(),
			Error:      "Failed to initialize Start.GG API client: " + err.Error(),
		}))
	}

	// Fetch user tournament history
	results, err := client.GetUserTournamentHistory("mckusa", 10) // Fetch last 10 tournaments
	if err != nil {
		h.Logger.Error("Failed to fetch tournament history", zap.Error(err))
		return Render(c, 200, panel.PanelStartGG(&panel.StartGGPanelData{
			UserSlug:   "mckusa",
			LastUpdate: time.Now(),
			Error:      "Failed to fetch tournament data: " + err.Error(),
		}))
	}

	// Get user info from the first result if available
	userName := ""
	if len(results) > 0 {
		// Note: The API structure might not include user name in this query
		// We could make a separate query for user info if needed
		userName = "mckusa" // Default for now
	}

	return Render(c, 200, panel.PanelStartGG(&panel.StartGGPanelData{
		UserSlug:   "mckusa",
		UserName:   userName,
		Results:    results,
		LastUpdate: time.Now(),
	}))
}
