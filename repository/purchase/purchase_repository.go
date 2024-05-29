package purchase_repository

import (
	"context"

	purchase_entity "github.com/kangman53/project-sprint-belibang/entity/Purchase"
)

type PurchaseRepository interface {
	Estimate(ctx context.Context, purchases purchase_entity.Purchase) (purchase_entity.Purchase, error)
	Order(ctx context.Context, req purchase_entity.Purchase) (purchase_entity.Purchase, error)
}
