package controller

import (
	"github.com/gofiber/fiber/v2"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	merchant_service "github.com/kangman53/project-sprint-belibang/service/merchant"
)

type MerchantController struct {
	MerchantService merchant_service.MerchantService
}

func NewMerchantController(merchantService merchant_service.MerchantService) *MerchantController {
	return &MerchantController{
		MerchantService: merchantService,
	}
}

func (controller MerchantController) Add(ctx *fiber.Ctx) error {
	merchantReq := new(merchant_entity.AddMerchantRequest)
	if err := ctx.BodyParser(merchantReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.MerchantService.Add(ctx, *merchantReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (controller MerchantController) SearchNearby(ctx *fiber.Ctx) error {
	searchQuery := new(merchant_entity.SearchNearbyMerchantQuery)
	searchQuery.Limit = 5
	searchQuery.Offset = 0

	if err := ctx.QueryParser(searchQuery); err != nil {
		return exc.BadRequestException("Error when parsing request query")
	}

	resp, err := controller.MerchantService.SearchNearby(ctx, *searchQuery)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
