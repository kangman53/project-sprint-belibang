package merchant_repository

import (
	"context"

	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type MerchantRepository interface {
	Add(ctx context.Context, req merchant_entity.Merchant) (string, error)
	Search(ctx context.Context, query merchant_entity.SearchMerchantQuery) (*[]merchant_entity.SearchMerchantData, error)
}
