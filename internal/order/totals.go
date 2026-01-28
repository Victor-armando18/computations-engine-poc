package order

import "github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"

func recomputeTotals(order *canonical.OrderState) {
	sub := 0.0

	for _, it := range order.Items {
		sub += float64(it.Quantity) * it.Price
	}

	order.Totals.SubTotal = sub
	order.Totals.Total = sub
}
