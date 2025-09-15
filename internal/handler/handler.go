package internal

import (
	"daily-driver/internal/api/art"
	"daily-driver/web/static/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handler struct {
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

func (h *Handler) AttachRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return Render(c, 200, templates.Index())
	})
	e.GET("/artwork", h.GetArtwork)
}

func (h *Handler) GetArtwork(c echo.Context) error {
	artAPI := art.NewChicagoAPIClient()

	hasImage := false
	artwork := &art.Artwork{}
	for hasImage == false {
		temp, err := artAPI.GetRandomArtwork()
		if err != nil {
			return c.String(500, "Error fetching artwork")
		}
		if artwork.ImageURL() != "" {
			hasImage = true
		}
		artwork = temp
	}

	return Render(c, 200, templates.ArtPanel(artwork))
}
