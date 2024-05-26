package controller

import (
	"github.com/gofiber/fiber/v2"
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	item_service "github.com/kangman53/project-sprint-belibang/service/item"
)

type ItemController struct {
	ItemService item_service.ItemService
}

func NewItemController(itemService item_service.ItemService) *ItemController {
	return &ItemController{
		ItemService: itemService,
	}
}

func (controller ItemController) Add(ctx *fiber.Ctx) error {
	itemReq := new(item_entity.AddItemRequest)
	if err := ctx.BodyParser(itemReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.ItemService.Add(ctx, *itemReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
