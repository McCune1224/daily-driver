package internal

import (
	"daily-driver/web/static/templates"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const PanelRotationIntervalSeconds = 3

func (h *Handler) RenderPanels(c echo.Context) error {
	// Define panel options
	panelOptions := []func(echo.Context) error{
		func(c echo.Context) error { return Render(c, 200, templates.Boilerplate(0)) },
		func(c echo.Context) error { return Render(c, 200, templates.Boilerplate(1)) },
		// func(c echo.Context) error { return h.GetRandomArtwork(c) },
	}

	// Determine the current panel index based on time
	numPanels := len(panelOptions) // Total number of panels available
	if numPanels == 0 {
		return c.String(500, "No panels available for rotation")
	}

	currentTime := time.Now().Unix()
	panelIndex := int((currentTime / PanelRotationIntervalSeconds) % int64(numPanels))

	h.Logger.Info("Rotating panel",
		zap.Int("currentPanelIndex", panelIndex),
	)

	return panelOptions[panelIndex](c)
}
