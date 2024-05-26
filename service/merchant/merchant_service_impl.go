package merchant_service

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	merchantRep "github.com/kangman53/project-sprint-belibang/repository/merchant"
)

type merchantServiceImpl struct {
	MerchantRepository merchantRep.MerchantRepository
	Validator          *validator.Validate
}

func NewMerchantService(merchantRepository merchantRep.MerchantRepository, validator *validator.Validate) MerchantService {
	return &merchantServiceImpl{
		MerchantRepository: merchantRepository,
		Validator:          validator,
	}
}

func (service *merchantServiceImpl) Add(ctx *fiber.Ctx, req merchant_entity.MerchantRegisterRequest) (merchant_entity.MerchantRegisterResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return merchant_entity.MerchantRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	merchant := merchant_entity.Merchant{
		Name:      req.Name,
		Category:  req.Category,
		ImageUrl:  req.ImageUrl,
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	userCtx := ctx.UserContext()
	merchantId, err := service.MerchantRepository.Add(userCtx, merchant)
	if err != nil {
		return merchant_entity.MerchantRegisterResponse{}, err
	}

	return merchant_entity.MerchantRegisterResponse{
		Id: merchantId,
	}, nil
}
