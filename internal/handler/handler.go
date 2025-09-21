package internal

import (
	"daily-driver/web/static/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// TODO: Add DB stuff...
type Handler struct {
	Logger *zap.Logger
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

	// Panel routes
	panel := e.Group("/panel")
	panel.GET("", h.RenderPanels)

	// Art routes
	art := e.Group("/art")
	art.GET("/api/random", h.GetRandomArtwork)
}
