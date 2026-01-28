package rules

import (
	"encoding/json"
	"fmt"

	"github.com/dolphin-sistemas/computations-engine/core"
)

func MapToCoreState(m map[string]interface{}) (core.State, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return core.State{}, err
	}

	var state core.State
	if err := json.Unmarshal(raw, &state); err != nil {
		return core.State{}, fmt.Errorf("cannot decode engine result to core.State: %w", err)
	}

	return state, nil
}
