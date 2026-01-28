package rules

import (
	"context"
	"errors"

	"github.com/Victor-armando18/computations-engine-poc/internal/order/canonical"
	engine "github.com/dolphin-sistemas/computations-engine"
	"github.com/dolphin-sistemas/computations-engine/core"
)

var ErrRuleViolation = errors.New("rule violation")

func Run(
	ctx context.Context,
	state canonical.OrderState,
	meta core.ContextMeta,
) (canonical.OrderState, []core.Reason, error) {

	// 1. Converter estado de domínio → core
	coreState := ToCoreState(state)

	// 2. Carregar RulePack (JSON)
	rulePack, err := LoadRulePack()
	if err != nil {
		return state, nil, err
	}

	// 3. Executar engine
	resultState,
		_,
		reasons,
		violations,
		_,
		err := engine.RunEngine(ctx, coreState, rulePack, meta)

	if err != nil {
		return state, nil, err
	}

	// 4. Se houve violations → bloquear
	if len(violations) > 0 {
		return state, reasons, ErrRuleViolation
	}

	// 5. Converter core → domínio
	coreFinal, err := MapToCoreState(resultState)
	if err != nil {
		return state, reasons, err
	}

	return FromCoreState(coreFinal), reasons, nil
}
