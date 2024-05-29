package purchase_service

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	purchase_entity "github.com/kangman53/project-sprint-belibang/entity/Purchase"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	purchase_repository "github.com/kangman53/project-sprint-belibang/repository/purchase"
	auth_service "github.com/kangman53/project-sprint-belibang/service/auth"
)

type purchaseServiceImpl struct {
	PurchaseRepository purchase_repository.PurchaseRepository
	AuthService        auth_service.AuthService
	Validator          *validator.Validate
}

func NewPurchaseService(purchaseRepository purchase_repository.PurchaseRepository, authService auth_service.AuthService, validator *validator.Validate) PurchaseService {
	return &purchaseServiceImpl{
		PurchaseRepository: purchaseRepository,
		AuthService:        authService,
		Validator:          validator,
	}
}

func (service *purchaseServiceImpl) Estimate(ctx *fiber.Ctx, req purchase_entity.PurchaseEstimateRequest) (purchase_entity.PurchaseEstimateResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return purchase_entity.PurchaseEstimateResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	userId, err := service.AuthService.GetValidUser(ctx)
	if err != nil {
		return purchase_entity.PurchaseEstimateResponse{}, exc.UnauthorizedException("Unauthorized")
	}

	purchase := purchase_entity.Purchase{
		UserId:    userId,
		Latitude:  req.UserLocation.Latitude,
		Longitude: req.UserLocation.Longitude,
		Order:     *req.Orders,
	}

	userContext := ctx.Context()
	purchaseInserted, err := service.PurchaseRepository.Estimate(userContext, purchase)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return purchase_entity.PurchaseEstimateResponse{}, exc.BadRequestException("distance too far")
		}
		if strings.Contains(err.Error(), "orders_merchant_id_fkey") {
			return purchase_entity.PurchaseEstimateResponse{}, exc.NotFoundException("merchantId not found")
		}
		if strings.Contains(err.Error(), "not found") {
			return purchase_entity.PurchaseEstimateResponse{}, exc.NotFoundException("itemId not found")
		}
		return purchase_entity.PurchaseEstimateResponse{}, err
	}

	return purchase_entity.PurchaseEstimateResponse{
		TotalPrice:                     purchaseInserted.TotalPrice,
		EstimatedDeliveryTimeInMinutes: purchaseInserted.EstimatedDeliveryTimeInMinutes,
		CalculatedEstimateId:           purchaseInserted.CalculatedEstimateId,
	}, nil
}

func (service *purchaseServiceImpl) Order(ctx *fiber.Ctx, req purchase_entity.PurchaseOrderRequest) (purchase_entity.PurchaseOrderResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return purchase_entity.PurchaseOrderResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	userId, err := service.AuthService.GetValidUser(ctx)
	if err != nil {
		return purchase_entity.PurchaseOrderResponse{}, exc.UnauthorizedException("Unauthorized")
	}

	purchase := purchase_entity.Purchase{
		UserId:               userId,
		CalculatedEstimateId: req.CalculatedEstimateId,
	}

	userContext := ctx.Context()
	purchaseUpdated, err := service.PurchaseRepository.Order(userContext, purchase)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return purchase_entity.PurchaseOrderResponse{}, exc.NotFoundException("calculatedEstimateId not found")
		}
	}

	return purchase_entity.PurchaseOrderResponse{
		OrderId: purchaseUpdated.Id,
	}, nil
}
