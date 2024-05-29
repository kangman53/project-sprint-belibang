package controller

import (
	"github.com/gofiber/fiber/v2"
	purchase_entity "github.com/kangman53/project-sprint-belibang/entity/Purchase"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	purchase_service "github.com/kangman53/project-sprint-belibang/service/purchase"
)

type PurchaseController struct {
	PurchaseService purchase_service.PurchaseService
}

func NewPurchaseController(purchaseService purchase_service.PurchaseService) *PurchaseController {
	return &PurchaseController{
		PurchaseService: purchaseService,
	}
}

func (controller PurchaseController) Estimate(ctx *fiber.Ctx) error {
	purchaseReq := new(purchase_entity.PurchaseEstimateRequest)
	if err := ctx.BodyParser(purchaseReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.PurchaseService.Estimate(ctx, *purchaseReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller PurchaseController) Order(ctx *fiber.Ctx) error {
	purchaseReq := new(purchase_entity.PurchaseOrderRequest)
	if err := ctx.BodyParser(purchaseReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.PurchaseService.Order(ctx, *purchaseReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
