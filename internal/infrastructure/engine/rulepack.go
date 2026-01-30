package engine

import (
	"encoding/json"
	"os"

	"github.com/dolphin-sistemas/computations-engine/core"
)

// LoadRulePack carrega o arquivo JSON das regras
func LoadRulePack(path string) (core.RulePack, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return core.RulePack{}, err
	}

	var pack core.RulePack
	if err := json.Unmarshal(raw, &pack); err != nil {
		return core.RulePack{}, err
	}

	return pack, nil
}
