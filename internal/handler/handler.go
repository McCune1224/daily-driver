package internal

import (
	"daily-driver/web/static/templates"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// TODO: Add DB stuff...
type Handler struct {
	Logger *zap.Logger
	DBPool *pgxpool.Pool
	// GarminData []*proto.FIT
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
	e.GET("/", h.RenderIndex)

	// Admin routes
	admin := e.Group("/admin")
	admin.GET("", h.RenderAdmin)
	admin.POST("/upload/garmin", h.UploadGarminFile)
	// admin.GET("/garmin-data", h.GetGarminData)

	// Panel routes
	panel := e.Group("/panel")
	panel.GET("", h.RenderPanels)

	// Art routes
	art := e.Group("/art")
	art.GET("/api/random", h.GetRandomArtwork)
}
