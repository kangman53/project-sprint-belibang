package item_service

import (
	"github.com/gofiber/fiber/v2"
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
)

type ItemService interface {
	Add(ctx *fiber.Ctx, req item_entity.AddItemRequest) (item_entity.AddItemResponse, error)
	Search(ctx *fiber.Ctx, req item_entity.SearchItemQuery) (item_entity.SearchItemResponse, error)
}
