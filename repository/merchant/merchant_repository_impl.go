package merchant_repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type merchantRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewMerchantRepository(dbPool *pgxpool.Pool) MerchantRepository {
	return &merchantRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *merchantRepositoryImpl) Add(ctx context.Context, merchant merchant_entity.Merchant) (string, error) {
	var merchantId string
	query := "INSERT INTO merchants (name, category, image_url, latitude, longitude) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, merchant.Name, merchant.Category, merchant.ImageUrl, merchant.Latitude, merchant.Longitude).Scan(&merchantId); err != nil {
		return "", err
	}

	return merchantId, nil
}
