package usecase

import (
	"context"

	"github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
)

type OrderUseCase struct {
	repo   order.Repository
	engine order.RuleEngine
}

func NewOrderUseCase(repo order.Repository, engine order.RuleEngine) *OrderUseCase {
	return &OrderUseCase{repo: repo, engine: engine}
}

// ExecutePatch realiza o fluxo completo: Busca -> Patch -> Engine -> Salva -> Retorna
func (uc *OrderUseCase) ExecutePatch(
	ctx context.Context,
	id string,
	patchFn func(order.Order) (order.Order, error),
	meta order.ExecutionMeta,
) (order.RulesResult, error) {

	o, err := uc.repo.Get(ctx, id)
	if err != nil {
		return order.RulesResult{}, err
	}

	patched, err := patchFn(o)
	if err != nil {
		return order.RulesResult{}, err
	}

	result, err := uc.engine.Run(ctx, patched, meta)
	if err != nil {
		return order.RulesResult{}, err
	}

	// Persistência (Exclusivo do Patch)
	if err := uc.repo.Save(ctx, result.StateFragment); err != nil {
		return order.RulesResult{}, err
	}

	return result, nil
}

// ExecuteValidate realiza o fluxo: Busca -> Patch -> Engine -> Retorna (Sem Salvar)
func (uc *OrderUseCase) ExecuteValidate(
	ctx context.Context,
	id string,
	patchFn func(order.Order) (order.Order, error),
	meta order.ExecutionMeta,
) (order.RulesResult, error) {

	o, err := uc.repo.Get(ctx, id)
	if err != nil {
		return order.RulesResult{}, err
	}

	patched, err := patchFn(o)
	if err != nil {
		return order.RulesResult{}, err
	}

	// Apenas simula o processamento das regras
	return uc.engine.Run(ctx, patched, meta)
}
