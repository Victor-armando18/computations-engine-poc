package usecase

import (
	"context"

	domain "github.com/Victor-armando18/computations-engine-poc/internal/domain/order"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/engine"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/repository"
)

type PatchOrder struct {
	repo   *repository.OrderMemoryRepository
	engine *engine.Adapter
}

func NewPatchOrder(repo *repository.OrderMemoryRepository, eng *engine.Adapter) *PatchOrder {
	return &PatchOrder{repo: repo, engine: eng}
}

func (uc *PatchOrder) Execute(
	ctx context.Context,
	id string,
	apply func(domain.Order) (domain.Order, error),
	meta domain.ExecutionMeta,
) (domain.RulesResult, error) {

	order, err := uc.repo.Get(id)
	if err != nil {
		return domain.RulesResult{}, err
	}

	patched, err := apply(order)
	if err != nil {
		return domain.RulesResult{}, err
	}

	result, err := uc.engine.Run(ctx, patched, meta)
	if err != nil {
		return domain.RulesResult{}, err
	}

	uc.repo.Save(result.StateFragment)

	return result, nil
}
