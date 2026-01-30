package http

import (
	"encoding/json"

	domain "github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	jsonpatch "github.com/evanphx/json-patch/v5"
)

func ApplyJSONPatch(o domain.Order, patch []byte) (domain.Order, error) {
	raw, err := json.Marshal(o)
	if err != nil {
		return domain.Order{}, err
	}

	p, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		return domain.Order{}, err
	}

	modified, err := p.Apply(raw)
	if err != nil {
		return domain.Order{}, err
	}

	var out domain.Order
	if err := json.Unmarshal(modified, &out); err != nil {
		return domain.Order{}, err
	}

	return out, nil
}
