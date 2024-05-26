package item_service

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	itemRep "github.com/kangman53/project-sprint-belibang/repository/item"
)

type itemServiceImpl struct {
	ItemRepository itemRep.ItemRepository
	Validator      *validator.Validate
}

func NewMerchantService(itemRepository itemRep.ItemRepository, validator *validator.Validate) ItemService {
	return &itemServiceImpl{
		ItemRepository: itemRepository,
		Validator:      validator,
	}
}

func (service *itemServiceImpl) Add(ctx *fiber.Ctx, req item_entity.AddItemRequest) (item_entity.AddItemResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return item_entity.AddItemResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	merchant := item_entity.Item{
		Name:     req.Name,
		Category: req.Category,
		ImageUrl: req.ImageUrl,
		Price:    req.Price,
	}

	userCtx := ctx.UserContext()
	itemId, err := service.ItemRepository.Add(userCtx, merchant)
	if err != nil {
		return item_entity.AddItemResponse{}, err
	}

	return item_entity.AddItemResponse{
		Id: itemId,
	}, nil
}
