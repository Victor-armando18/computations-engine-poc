package order

import (
	"context"
	"fmt"

	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	"github.com/Victor-armando18/computations-engine-poc/internal/order/idempotency"
	"github.com/Victor-armando18/computations-engine-poc/internal/order/patch"
	"github.com/Victor-armando18/computations-engine-poc/internal/order/rules"
	"github.com/Victor-armando18/computations-engine-poc/internal/response"
	"github.com/dolphin-sistemas/computations-engine/core"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

type Service struct {
	store *idempotency.Store[response.Envelope]
}

func NewService() *Service {
	return &Service{
		store: idempotency.New[response.Envelope](),
	}
}

// Patch aplica um JSON Patch no estado da ordem e retorna Envelope
func (s *Service) Patch(ctx context.Context, key string, state canonical.OrderState, body []byte) (response.Envelope, error) {
	// verificar cache idempotency
	if cached, ok := s.store.Get(key); ok {
		return cached, nil
	}

	// decodificar array de operações
	p, err := jsonpatch.DecodePatch(body)
	if err != nil {
		return response.Envelope{}, fmt.Errorf("invalid patch body: %w", err)
	}

	// validar tamanho
	if err := patch.ValidatePatchSize(len(p)); err != nil {
		return response.Envelope{}, err
	}

	// garantir array inicializado
	if state.Items == nil {
		state.Items = []canonical.OrderItem{}
	}

	// criar placeholders para índices inexistentes se patch referenciar
	for _, op := range p {
		path, err := op.Path()
		if err != nil {
			continue
		}

		var idx int
		n, _ := fmt.Sscanf(path, "/items/%d", &idx)
		if n == 1 {
			for len(state.Items) <= idx {
				state.Items = append(state.Items, canonical.OrderItem{})
			}
		}
	}

	// aplicar patch
	draft, err := patch.Apply(state, p)
	if err != nil {
		return response.Envelope{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	draft = patch.RecalculateTotals(draft)
	// executar regras
	final, reasons, err := rules.Run(ctx, draft, core.ContextMeta{})
	if err != nil {
		return response.Envelope{}, err
	}

	// construir Envelope
	resp := response.Build(final, reasons)

	// salvar cache idempotency
	s.store.Set(key, resp)

	return resp, nil
}

// Validate apenas executa regras sem alterar estado
func (s *Service) Validate(ctx context.Context, state canonical.OrderState, _ []byte) (response.Envelope, error) {
	final, reasons, err := rules.Run(ctx, state, core.ContextMeta{})
	if err != nil {
		return response.Envelope{}, err
	}

	resp := response.Build(final, reasons)
	return resp, nil
}
