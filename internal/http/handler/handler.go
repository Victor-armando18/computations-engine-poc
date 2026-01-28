package handler

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/Victor-armando18/computations-engine-poc/internal/order"
	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *order.Service
	mockDB  map[string]canonical.OrderState
	mu      sync.RWMutex
}

func New() *Handler {
	return &Handler{
		service: order.NewService(),
		mockDB: map[string]canonical.OrderState{
			"123": {
				ID: "123",
				Items: []canonical.OrderItem{
					{SKU: "ITEM-001", Quantity: 2, Price: 50},
					{SKU: "ITEM-002", Quantity: 1, Price: 100},
				},
				Totals: canonical.Totals{
					SubTotal: 2*50 + 1*100,
					Total:    2*50 + 1*100,
				},
			},
		},
	}
}

// PATCH endpoint: aplica JSON Patch na ordem
func (h *Handler) Patch(c echo.Context) error {
	key := c.Request().Header.Get("Idempotency-Key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Idempotency-Key header required"})
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.mu.RLock()
	state, ok := h.mockDB[c.Param("id")]
	h.mu.RUnlock()
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	resp, err := h.service.Patch(context.Background(), key, state, body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.mu.Lock()
	h.mockDB[c.Param("id")] = resp.StateFragment
	h.mu.Unlock()

	return c.JSON(http.StatusOK, resp)
}

// VALIDATE endpoint: apenas roda regras sem alterar estado
func (h *Handler) Validate(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.mu.RLock()
	state, ok := h.mockDB[c.Param("id")]
	h.mu.RUnlock()
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	resp, err := h.service.Validate(context.Background(), state, body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
