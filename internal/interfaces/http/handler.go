package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	domain "github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/labstack/echo/v4"
)

type RFC7807Error struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

type OrderHandler struct {
	Patch interface {
		Execute(context.Context, string, func(domain.Order) (domain.Order, error), domain.ExecutionMeta) (domain.RulesResult, error)
	}
	Validate interface {
		Execute(context.Context, string, func(domain.Order) (domain.Order, error), domain.ExecutionMeta) (domain.RulesResult, error)
	}
}

func (h *OrderHandler) PatchOrder(c echo.Context) error {
	return h.handle(c, h.Patch.Execute)
}

func (h *OrderHandler) ValidateOrder(c echo.Context) error {
	return h.handle(c, h.Validate.Execute)
}

func (h *OrderHandler) handle(c echo.Context, exec func(context.Context, string, func(domain.Order) (domain.Order, error), domain.ExecutionMeta) (domain.RulesResult, error)) error {
	// Timeout total da requisição HTTP (incluindo retries do motor)
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	// Regra de segurança: limite de tamanho do patch
	body, err := io.ReadAll(io.LimitReader(c.Request().Body, 1024*1024))
	if err != nil {
		return h.sendError(c, http.StatusBadRequest, "Payload Error", "Request body too large or unreadable")
	}

	meta := domain.ExecutionMeta{
		ContextGuid:    c.Request().Header.Get("ContextGuid"),
		UserCode:       c.Request().Header.Get("userCode"),
		CorrelationID:  c.Request().Header.Get("correlation-id"),
		IdempotencyKey: c.Request().Header.Get("Idempotency-Key"),
		Locale:         c.Request().Header.Get("Accept-Language"),
	}

	apply := func(o domain.Order) (domain.Order, error) {
		raw, _ := json.Marshal(o)
		patch, err := jsonpatch.DecodePatch(body)
		if err != nil {
			return o, domain.ErrInvalidPatch
		}

		mod, err := patch.Apply(raw)
		if err != nil {
			return o, err
		}

		var out domain.Order
		json.Unmarshal(mod, &out)
		return out, nil
	}

	res, err := exec(ctx, id, apply, meta)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			return h.sendError(c, http.StatusNotFound, "Order Not Found", err.Error())
		}
		return h.sendError(c, http.StatusInternalServerError, "Execution Failed", err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *OrderHandler) sendError(c echo.Context, status int, title, detail string) error {
	return c.JSON(status, RFC7807Error{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Path(),
	})
}
