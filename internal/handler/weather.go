package handler

import (
	"daily-driver/internal/api/weather"
	"daily-driver/web/static/templates/panel"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RenderPanelWeather renders the visual weather panel (HTML via templ)
func (h *Handler) RenderPanelWeather(c echo.Context) error {
	client := weather.NewClient()

	resp, err := client.GetRochesterWeather(c.Request().Context())
	if err != nil {
		if h.Logger != nil {
			h.Logger.Error("failed to fetch weather", zap.Error(err))
		}
		return c.JSON(502, map[string]string{"error": "failed to fetch weather"})
	}
	return Render(c, 200, panel.PanelWeather(resp))
}
