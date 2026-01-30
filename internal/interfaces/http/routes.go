package http

import "github.com/labstack/echo/v4"

func Register(e *echo.Echo, h *OrderHandler) {
	e.POST("/orders/:id/patch", func(c echo.Context) error {
		return h.handle(c, h.Patch.Execute)
	})
	e.POST("/orders/:id/validate", func(c echo.Context) error {
		return h.handle(c, h.Validate.Execute)
	})
}
