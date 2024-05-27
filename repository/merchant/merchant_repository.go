package merchant_repository

import (
	"context"

	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type MerchantRepository interface {
	Add(ctx context.Context, req merchant_entity.Merchant) (string, error)
	SearchNearby(ctx context.Context, query merchant_entity.SearchNearbyMerchantQuery, latitude float64, longitude float64) (*[]merchant_entity.SearchNearbyMerchantData, error)
}
