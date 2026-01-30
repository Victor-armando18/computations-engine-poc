package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Victor-armando18/computations-engine-poc/internal/application/usecase"
	"github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	useCase *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: uc}
}

func (h *OrderHandler) HandlePatch(c echo.Context) error {
	return h.process(c, h.useCase.ExecutePatch)
}

func (h *OrderHandler) HandleValidate(c echo.Context) error {
	return h.process(c, h.useCase.ExecuteValidate)
}

// process é um helper privado para evitar duplicação de lógica de extração de headers e body (DRY)
func (h *OrderHandler) process(c echo.Context, action func(context.Context, string, func(order.Order) (order.Order, error), order.ExecutionMeta) (order.RulesResult, error)) error {
	id := c.Param("id")

	// Segurança: Limite de Payload
	body, err := io.ReadAll(io.LimitReader(c.Request().Body, 1024*1024))
	if err != nil {
		return h.sendProblem(c, http.StatusBadRequest, "Payload Error", "Unable to read body")
	}

	meta := order.ExecutionMeta{
		ContextGuid:    c.Request().Header.Get("ContextGuid"),
		UserCode:       c.Request().Header.Get("userCode"),
		CorrelationID:  c.Request().Header.Get("correlation-id"),
		IdempotencyKey: c.Request().Header.Get("Idempotency-Key"),
		Locale:         c.Request().Header.Get("Accept-Language"),
	}

	// Closure para aplicar o patch RFC 6902
	patcher := func(o order.Order) (order.Order, error) {
		rawOrder, _ := json.Marshal(o)
		patch, err := jsonpatch.DecodePatch(body)
		if err != nil {
			return o, order.ErrInvalidPatch
		}
		modified, err := patch.Apply(rawOrder)
		if err != nil {
			return o, err
		}
		var updated order.Order
		json.Unmarshal(modified, &updated)
		return updated, nil
	}

	result, err := action(c.Request().Context(), id, patcher, meta)
	if err != nil {
		if err == order.ErrOrderNotFound {
			return h.sendProblem(c, http.StatusNotFound, "Not Found", err.Error())
		}
		return h.sendProblem(c, http.StatusInternalServerError, "Execution Error", err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) sendProblem(c echo.Context, status int, title, detail string) error {
	return c.JSON(status, map[string]any{
		"type":     "about:blank",
		"title":    title,
		"status":   status,
		"detail":   detail,
		"instance": c.Path(),
	})
}
