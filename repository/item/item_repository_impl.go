package item_repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
)

type itemRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewItemRepository(dbPool *pgxpool.Pool) ItemRepository {
	return &itemRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *itemRepositoryImpl) Add(ctx context.Context, item item_entity.Item, merchantId string) (string, error) {
	var itemId string
	query := "INSERT INTO items (name, category, image_url, price, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, item.Name, item.Category, item.ImageUrl, item.Price, merchantId).Scan(&itemId); err != nil {
		return "", err
	}

	return itemId, nil
}
