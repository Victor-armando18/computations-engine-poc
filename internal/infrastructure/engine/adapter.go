package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	domain "github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	engine "github.com/dolphin-sistemas/computations-engine"
	"github.com/dolphin-sistemas/computations-engine/core"
)

type Adapter struct {
	rulePack core.RulePack
	cache    sync.Map
}

func NewAdapter(rulePath string) (*Adapter, error) {
	pack, err := LoadRulePack(rulePath)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		rulePack: pack,
		cache:    sync.Map{},
	}, nil
}

func (a *Adapter) Run(
	ctx context.Context,
	order domain.Order,
	meta domain.ExecutionMeta,
) (domain.RulesResult, error) {

	// 1. Idempotência via Idempotency-Key
	if meta.IdempotencyKey != "" {
		if val, ok := a.cache.Load(meta.IdempotencyKey); ok {
			return val.(domain.RulesResult), nil
		}
	}

	base := toCoreState(order)

	// 2. Configuração de Retries e Timeouts
	var result *core.RunEngineResult
	var err error
	maxRetries := 3
	backoff := 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		// Criar um sub-contexto com timeout para cada tentativa (ex: 2 segundos)
		runCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

		// Chamada para a lib conforme a assinatura fornecida
		result, err = engine.RunEngine(
			runCtx,
			base,
			a.rulePack,
			core.ContextMeta{
				TenantID: meta.ContextGuid,
				UserID:   meta.UserCode,
				Locale:   meta.Locale,
			},
		)
		cancel()

		if err == nil {
			break
		}

		// Se for a última tentativa e deu erro, retorna
		if i == maxRetries-1 {
			return domain.RulesResult{}, fmt.Errorf("engine failed after %d attempts: %w", maxRetries, err)
		}

		// Espera antes da próxima tentativa
		select {
		case <-ctx.Done():
			return domain.RulesResult{}, ctx.Err()
		case <-time.After(backoff):
			backoff *= 2 // Exponential backoff
		}
	}

	// 3. Processamento do Resultado
	// Aplicar o StateFragment (map[string]interface{}) sobre o estado base
	finalCore, err := applyFragment(base, result.StateFragment)
	if err != nil {
		return domain.RulesResult{}, err
	}

	finalOrder := fromCoreState(finalCore)

	// Usar o ServerDelta retornado pelo motor se disponível,
	// ou calcular manualmente se necessário. Aqui usamos o do motor.
	response := domain.RulesResult{
		StateFragment: finalOrder,
		ServerDelta:   result.ServerDelta,
		Reasons:       result.Reasons,
		Violations:    result.Violations,
		RulesVersion:  a.rulePack.Version,
		CorrelationID: meta.CorrelationID,
	}

	// 4. Cache para Idempotência
	if meta.IdempotencyKey != "" {
		a.cache.Store(meta.IdempotencyKey, response)
	}

	return response, nil
}

func toCoreState(o domain.Order) core.State {
	items := make([]core.Item, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, core.Item{
			ID:     it.SKU,
			Amount: float64(it.Quantity),
			Fields: map[string]any{
				"sku":      it.SKU,
				"quantity": it.Quantity,
				"price":    it.Price,
			},
		})
	}

	return core.State{
		ID:    o.ID,
		Items: items,
		Totals: core.Totals{
			Subtotal: o.Totals.SubTotal,
			Total:    o.Totals.Total,
		},
	}
}

func fromCoreState(s core.State) domain.Order {
	items := make([]domain.OrderItem, len(s.Items))
	for i, it := range s.Items {
		price, _ := it.Fields["price"].(float64)
		qty := int(it.Amount)
		if q, ok := it.Fields["quantity"].(float64); ok && qty == 0 {
			qty = int(q)
		}

		items[i] = domain.OrderItem{
			SKU:      it.ID,
			Quantity: qty,
			Price:    price,
		}
	}

	return domain.Order{
		ID:    s.ID,
		Items: items,
		Totals: domain.Totals{
			SubTotal: s.Totals.Subtotal,
			Total:    s.Totals.Total,
		},
	}
}

func applyFragment(base core.State, frag map[string]interface{}) (core.State, error) {
	// Se o fragmento for nulo, retorna a base
	if frag == nil {
		return base, nil
	}

	b, _ := json.Marshal(base)
	var merged map[string]interface{}
	json.Unmarshal(b, &merged)

	// Merge dos campos alterados
	for k, v := range frag {
		merged[k] = v
	}

	out, _ := json.Marshal(merged)
	var final core.State
	err := json.Unmarshal(out, &final)
	return final, err
}
