package repository

import (
	"sync"

	"github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
)

type OrderMemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]order.Order
}

func NewOrderMemoryRepository() *OrderMemoryRepository {
	return &OrderMemoryRepository{
		orders: map[string]order.Order{
			"123": {
				ID: "123",
				Items: []order.OrderItem{
					{SKU: "ITEM-001", Quantity: 2, Price: 50},
					{SKU: "ITEM-002", Quantity: 1, Price: 100},
				},
				Totals: order.Totals{
					SubTotal: 200,
					Total:    200,
				},
			},
		},
	}
}

func (r *OrderMemoryRepository) Get(id string) (order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.orders[id]
	if !ok {
		return order.Order{}, order.ErrOrderNotFound
	}
	return o, nil
}

func (r *OrderMemoryRepository) Save(o order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID] = o
	return nil
}
