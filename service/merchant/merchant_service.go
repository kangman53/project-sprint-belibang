package merchant_service

import (
	"github.com/gofiber/fiber/v2"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type MerchantService interface {
	Add(ctx *fiber.Ctx, req merchant_entity.MerchantRegisterRequest) (merchant_entity.MerchantRegisterResponse, error)
}
