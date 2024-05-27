package merchant_service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type MerchantService interface {
	Add(ctx *fiber.Ctx, req merchant_entity.AddMerchantRequest) (merchant_entity.AddMerchantResponse, error)
	Search(ctx context.Context, req merchant_entity.SearchMerchantQuery) (merchant_entity.SearchMerchantResponse, error)
}
