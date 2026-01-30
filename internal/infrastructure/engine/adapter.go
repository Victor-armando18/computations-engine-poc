package engine

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	engine "github.com/dolphin-sistemas/computations-engine"
	"github.com/dolphin-sistemas/computations-engine/core"
)

type Adapter struct {
	rulePack core.RulePack
	cache    sync.Map
}

func NewAdapter(rulePath string) (*Adapter, error) {
	raw, err := os.ReadFile(rulePath)
	if err != nil {
		return nil, err
	}

	var pack core.RulePack
	if err := json.Unmarshal(raw, &pack); err != nil {
		return nil, err
	}

	return &Adapter{rulePack: pack}, nil
}

func (a *Adapter) Run(ctx context.Context, o order.Order, meta order.ExecutionMeta) (order.RulesResult, error) {
	// Early Return: Idempotência
	if meta.IdempotencyKey != "" {
		if val, ok := a.cache.Load(meta.IdempotencyKey); ok {
			return val.(order.RulesResult), nil
		}
	}

	base := a.toCoreState(o)

	// Execução com Retries
	res, err := a.runWithRetry(ctx, base, meta)
	if err != nil {
		return order.RulesResult{}, err
	}

	// Transformação de saída
	finalCore, _ := a.applyFragment(base, res.StateFragment)
	finalOrder := a.fromCoreState(finalCore)

	result := order.RulesResult{
		StateFragment: finalOrder,
		ServerDelta:   res.ServerDelta,
		Reasons:       res.Reasons,
		Violations:    res.Violations,
		RulesVersion:  a.rulePack.Version,
		CorrelationID: meta.CorrelationID,
	}

	if meta.IdempotencyKey != "" {
		a.cache.Store(meta.IdempotencyKey, result)
	}

	return result, nil
}

func (a *Adapter) runWithRetry(ctx context.Context, state core.State, meta order.ExecutionMeta) (*core.RunEngineResult, error) {
	var res *core.RunEngineResult
	var err error

	backoff := 100 * time.Millisecond
	for i := 0; i < 3; i++ {
		timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		res, err = engine.RunEngine(timeoutCtx, state, a.rulePack, core.ContextMeta{
			TenantID: meta.ContextGuid,
			UserID:   meta.UserCode,
			Locale:   meta.Locale,
		})
		cancel()

		if err == nil {
			return res, nil
		}

		time.Sleep(backoff)
		backoff *= 2
	}
	return nil, err
}

func (a *Adapter) toCoreState(o order.Order) core.State {
	items := make([]core.Item, len(o.Items))
	for i, it := range o.Items {
		items[i] = core.Item{
			ID:     it.SKU,
			Amount: float64(it.Quantity),
			Fields: map[string]any{"sku": it.SKU, "quantity": it.Quantity, "price": it.Price},
		}
	}
	return core.State{
		ID: o.ID, Items: items,
		Totals: core.Totals{Subtotal: o.Totals.SubTotal, Total: o.Totals.Total},
	}
}

func (a *Adapter) fromCoreState(s core.State) order.Order {
	items := make([]order.OrderItem, len(s.Items))
	for i, it := range s.Items {
		price, _ := it.Fields["price"].(float64)
		qty := int(it.Amount)
		if q, ok := it.Fields["quantity"].(float64); ok && qty == 0 {
			qty = int(q)
		}
		items[i] = order.OrderItem{SKU: it.ID, Quantity: qty, Price: price}
	}
	return order.Order{
		ID: s.ID, Items: items,
		Totals: order.Totals{SubTotal: s.Totals.Subtotal, Total: s.Totals.Total},
	}
}

func (a *Adapter) applyFragment(base core.State, frag map[string]any) (core.State, error) {
	if frag == nil {
		return base, nil
	}
	b, _ := json.Marshal(base)
	var m map[string]any
	json.Unmarshal(b, &m)
	for k, v := range frag {
		m[k] = v
	}
	out, _ := json.Marshal(m)
	var final core.State
	json.Unmarshal(out, &final)
	return final, nil
}
