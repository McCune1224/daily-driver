package handler

import (
	"daily-driver/internal/routes"
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
	e.GET(routes.Root, h.RenderIndex)

	garmin := e.Group(routes.GarminBase)
	garmin.GET("", h.RenderPanelGarmin)
	garmin.POST(routes.GarminUpload, h.UploadGarminFile)

	// Panel routes
	panel := e.Group(routes.PanelBase)
	panel.GET("", h.RenderPanels)

	// Art routes
	art := e.Group(routes.ArtBase)
	art.GET(routes.ArtRandomAPI, h.RenderPanelArtwork)
}
