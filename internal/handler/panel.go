package handler

import (
	"daily-driver/internal/routes"
	"daily-driver/web/static/templates/panel"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) MapPanelRoutes(e *echo.Echo) {

}

func (h *Handler) RenderPanels(c echo.Context) error {
	// Define panel options
	numPanels := len(h.PanelHandlers) // Total number of panels available
	if numPanels == 0 {
		return c.String(500, "No panels available for rotation")
	}

	currentTime := time.Now().Unix()
	panelIndex := int((currentTime / int64(routes.PanelRotationIntervalSeconds)) % int64(numPanels))

	h.Logger.Info("Rotating panel",
		zap.Int("currentPanelIndex", panelIndex),
	)

	// //make quick POST request to update the panel index cookie
	// client := &http.Client{}
	// req, err := http.NewRequest("POST", routes.PanelIndex, nil)
	// if err != nil {
	// 	h.Logger.Error("Error creating request", zap.Error(err))
	// 	return c.String(500, "Internal Server Error")
	// }
	// _, err = client.Do(req)
	// h.UpdatePanelIndex(c)

	return h.PanelHandlers[panelIndex](c)
}

func (h *Handler) UpdatePanelIndex(c echo.Context) error {
	cookie, err := c.Cookie("panel_index")
	if err != nil {
		if err == echo.ErrCookieNotFound {
			// If the cookie is not found, set it to 0
			c.SetCookie(&http.Cookie{
				Name:  "panel_index",
				Value: "0",
			})
			return Render(c, 200, panel.PanelIndexDisplay(len(h.PanelHandlers), 0))
		}
		h.Logger.Error("Error retrieving panel_index cookie", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}

	if cookie == nil {
		h.Logger.Error("No cookie of panel_index", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}

	cVal := cookie.Value
	cIntVal, err := strconv.Atoi(cVal)
	if err != nil {
		h.Logger.Error("Error converting panel_index cookie to int", zap.Error(err))
		cIntVal = 0 // Default to 0 if conversion fails
	}
	cIntVal = (cIntVal + 1) % len(h.PanelHandlers)

	c.SetCookie(&http.Cookie{
		Name:  "panel_index",
		Value: strconv.Itoa(cIntVal),
	})
	return Render(c, 200, panel.PanelIndexDisplay(len(h.PanelHandlers), cIntVal))
}
