package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, h *OrderHandler) {
	// Endpoints exigidos na Task
	e.POST("/orders/:id/patch", h.HandlePatch)
	e.POST("/orders/:id/validate", h.HandleValidate)
}
