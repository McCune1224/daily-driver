package handler

import (
	"daily-driver/internal/api/art"
	"daily-driver/web/static/templates"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RenderPanelArtwork(c echo.Context) error {
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

	return Render(c, 200, templates.ComponentArtpiece(artwork))
}
