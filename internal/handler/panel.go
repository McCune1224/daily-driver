package handler

import (
	"daily-driver/internal/routes"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) RenderPanels(c echo.Context) error {
	// Define panel options
	panelOptions := []func(echo.Context) error{
		// func(c echo.Context) error { return Render(c, 200, templates.Boilerplate(0)) },
		func(c echo.Context) error { return h.RenderPanelWeather(c) },
		// func(c echo.Context) error { return h.RenderGarminPanel(c) },
		// func(c echo.Context) error { return h.RenderPanelArtwork(c) },
		// func(c echo.Context) error { return Render(c, 200, panel.PanelWeather()) },
	}

	numPanels := len(panelOptions) // Total number of panels available
	if numPanels == 0 {
		return c.String(500, "No panels available for rotation")
	}

	currentTime := time.Now().Unix()
	panelIndex := int((currentTime / int64(routes.PanelRotationIntervalSeconds)) % int64(numPanels))

	h.Logger.Info("Rotating panel",
		zap.Int("currentPanelIndex", panelIndex),
	)

	return panelOptions[panelIndex](c)
}
