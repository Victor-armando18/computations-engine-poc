package patch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

func ValidatePatchSize(size int) error {
	if size > 10 {
		return errors.New("patch too large")
	}
	return nil
}

func Apply(state canonical.OrderState, p jsonpatch.Patch) (canonical.OrderState, error) {
	if state.Items == nil {
		state.Items = []canonical.OrderItem{}
	}

	// criar placeholders para índices referenciados no patch
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

	raw, _ := json.Marshal(state)
	modified, err := p.Apply(raw)
	if err != nil {
		return state, err
	}

	var newState canonical.OrderState
	if err := json.Unmarshal(modified, &newState); err != nil {
		return state, err
	}

	if newState.Items == nil {
		newState.Items = []canonical.OrderItem{}
	}

	return newState, nil
}

func RecalculateTotals(state canonical.OrderState) canonical.OrderState {
	subTotal := 0.0
	for _, item := range state.Items {
		subTotal += float64(item.Quantity) * item.Price
	}
	state.Totals.SubTotal = subTotal
	state.Totals.Total = subTotal // se não houver outros descontos
	return state
}
