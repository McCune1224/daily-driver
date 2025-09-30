package handler

import (
	"daily-driver/web/static/templates/panel"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RenderPanelWeather(c echo.Context) error {
	return Render(c, 200, panel.PanelWeather())
}
