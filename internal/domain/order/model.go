package order

import "github.com/dolphin-sistemas/computations-engine/core"

type Order struct {
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

type ExecutionMeta struct {
	ContextGuid    string
	UserCode       string
	CorrelationID  string
	IdempotencyKey string
	Locale         string
}

type RulesResult struct {
	StateFragment Order            `json:"stateFragment"`
	ServerDelta   map[string]any   `json:"serverDelta"`
	Reasons       []core.Reason    `json:"reasons"`
	Violations    []core.Violation `json:"violations"`
	RulesVersion  string           `json:"rulesVersion"`
	CorrelationID string           `json:"correlationId"`
}
