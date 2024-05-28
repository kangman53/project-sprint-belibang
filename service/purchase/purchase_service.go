package purchase_service

import (
	"github.com/gofiber/fiber/v2"
	purchase_entity "github.com/kangman53/project-sprint-belibang/entity/Purchase"
)

type PurchaseService interface {
	Estimate(ctx *fiber.Ctx, req purchase_entity.PurchaseEstimateRequest) (purchase_entity.PurchaseEstimateResponse, error)
}
