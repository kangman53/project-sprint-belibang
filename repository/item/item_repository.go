package item_repository

import (
	"context"

	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
)

type ItemRepository interface {
	Add(ctx context.Context, req item_entity.Item, merchantId string) (string, error)
}
