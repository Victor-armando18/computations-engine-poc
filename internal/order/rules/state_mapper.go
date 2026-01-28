package rules

import (
	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	"github.com/dolphin-sistemas/computations-engine/core"
)

func ToCoreState(order canonical.OrderState) core.State {
	items := make([]core.Item, 0, len(order.Items))

	for _, it := range order.Items {
		items = append(items, core.Item{
			ID: it.SKU,
			Fields: map[string]interface{}{
				"quantity": float64(it.Quantity),
				"price":    it.Price,
			},
		})
	}

	return core.State{
		ID:    order.ID,
		Items: items,
		Totals: core.Totals{
			Subtotal: order.Totals.SubTotal,
			Total:    order.Totals.Total,
		},
		Fields: map[string]interface{}{
			"id": order.ID,
		},
	}
}

func FromCoreState(state core.State) canonical.OrderState {
	order := canonical.OrderState{
		Items:  []canonical.OrderItem{},
		Totals: canonical.Totals{},
	}

	if id, ok := state.Fields["id"].(string); ok {
		order.ID = id
	}

	for _, it := range state.Items {
		qty := 0
		if v, ok := it.Fields["quantity"].(float64); ok {
			qty = int(v)
		}

		price := 0.0
		if v, ok := it.Fields["price"].(float64); ok {
			price = v
		}

		order.Items = append(order.Items, canonical.OrderItem{
			SKU:      it.ID,
			Quantity: qty,
			Price:    price,
		})
	}

	order.Totals.SubTotal = state.Totals.Subtotal
	order.Totals.Total = state.Totals.Total

	return order
}
