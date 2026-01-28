package rules

import (
	"encoding/json"
	"os"

	"github.com/dolphin-sistemas/computations-engine/core"
)

func LoadRulePack() (core.RulePack, error) {
	data, err := os.ReadFile("rules.json")
	if err != nil {
		return core.RulePack{}, err
	}
	var pack core.RulePack
	if err := json.Unmarshal(data, &pack); err != nil {
		return core.RulePack{}, err
	}
	return pack, nil
}
