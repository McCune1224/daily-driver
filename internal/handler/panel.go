package handler

import (
	"daily-driver/internal/routes"
	"daily-driver/web/static/templates/panel"
	"net/http"
	"strconv"
	"strings"
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

	// Check panel_rotation cookie
	panelRotationCookie, err := c.Cookie("panel_rotation")
	rotationEnabled := true
	if err == nil && panelRotationCookie.Value == "off" {
		rotationEnabled = false
	}

	var panelIndex int
	if rotationEnabled {
		currentTime := time.Now().Unix()
		panelIndex = int((currentTime / int64(routes.PanelRotationIntervalSeconds)) % int64(numPanels))
	} else {
		// Use panel_index cookie
		panelIndexCookie, err := c.Cookie("panel_index")
		if err != nil || panelIndexCookie.Value == "" {
			// Default to 0 if not set
			c.SetCookie(&http.Cookie{
				Name:  "panel_index",
				Value: "0",
				Path:  "/",
			})
			panelIndex = 0
		} else {
			cVal := panelIndexCookie.Value
			cIntVal, err := strconv.Atoi(cVal)
			if err != nil {
				h.Logger.Error("Error converting panel_index cookie to int", zap.Error(err))
				cIntVal = 0
			}
			panelIndex = cIntVal % numPanels
		}
	}

	h.Logger.Info("Rendering panel",
		zap.Int("panelIndex", panelIndex),
		zap.Bool("rotationEnabled", rotationEnabled),
	)

	return h.PanelHandlers[panelIndex](c)
}

func (h *Handler) UpdatePanelIndex(c echo.Context) error {

	panelIndexCookie, err := c.Cookie("panel_index")
	if err != nil {
		if strings.Contains(err.Error(), "named cookie not present") {
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

	panelRotationCookie, err := c.Cookie("panel_rotation")
	if err != nil {
		if strings.Contains(err.Error(), "named cookie not present") {
			// If the cookie is not found, set it to "on"
			c.SetCookie(&http.Cookie{
				Name:  "panel_rotation",
				Value: "on",
				Path:  "/",
			})
			return Render(c, 200, panel.PanelIndexDisplay(len(h.PanelHandlers), 0))
		}
		h.Logger.Error("Error retrieving panel_rotation cookie", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}

	if panelRotationCookie.Value == "off" {
		// If rotation is off, do not update the index
		cVal := panelIndexCookie.Value
		cIntVal, err := strconv.Atoi(cVal)
		if err != nil {
			h.Logger.Error("Error converting panel_index cookie to int", zap.Error(err))
			cIntVal = 0 // Default to 0 if conversion fails
		}
		return Render(c, 200, panel.PanelIndexDisplay(len(h.PanelHandlers), cIntVal))
	}

	if panelIndexCookie == nil {
		h.Logger.Error("No cookie of panel_index", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}

	cVal := panelIndexCookie.Value
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

func (h *Handler) TogglePanelRotation(c echo.Context) error {
	cookie, err := c.Cookie("panel_rotation")
	if err != nil {
		if strings.Contains(err.Error(), "named cookie not present") {
			// If the cookie is not found, set it to "on"
			c.SetCookie(&http.Cookie{
				Name:  "panel_rotation",
				Path:  "/",
				Value: "on",
			})
			return c.String(200, "Panel rotation enabled")
		}
		h.Logger.Error("Error retrieving panel_rotation cookie", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}

	if cookie == nil {
		h.Logger.Error("No cookie of panel_rotation", zap.Error(err))
		return c.String(500, "Internal Server Error")
	}
	if cookie.Value == "on" {
		// Disable rotation
		c.SetCookie(&http.Cookie{
			Name:  "panel_rotation",
			Value: "off",
			Path:  "/",
		})
		return c.String(200, "Panel rotation disabled")
	} else {
		// Enable rotation
		c.SetCookie(&http.Cookie{
			Name:  "panel_rotation",
			Value: "on",
			Path:  "/",
		})
		return c.String(200, "Panel rotation enabled")
	}
}
