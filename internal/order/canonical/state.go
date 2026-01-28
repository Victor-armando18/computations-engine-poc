package canonical

type OrderState struct {
	ID     string      `json:"id"`
	Items  []OrderItem `json:"items"`
	Totals Totals      `json:"totals"`
}

type OrderItem struct {
	SKU      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Totals struct {
	SubTotal float64 `json:"subTotal"`
	Total    float64 `json:"total"`
}
