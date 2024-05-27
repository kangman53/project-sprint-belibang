package merchant_service

import (
	"context"
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

func (service *merchantServiceImpl) Add(ctx *fiber.Ctx, req merchant_entity.AddMerchantRequest) (merchant_entity.AddMerchantResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return merchant_entity.AddMerchantResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
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
		return merchant_entity.AddMerchantResponse{}, err
	}

	return merchant_entity.AddMerchantResponse{
		Id: merchantId,
	}, nil
}

func (service *merchantServiceImpl) Search(ctx context.Context, searchQuery merchant_entity.SearchMerchantQuery) (merchant_entity.SearchMerchantResponse, error) {
	if err := service.Validator.Struct(searchQuery); err != nil {
		return merchant_entity.SearchMerchantResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	merchantSearched, err := service.MerchantRepository.Search(ctx, searchQuery)
	if err != nil {
		return merchant_entity.SearchMerchantResponse{}, err
	}

	return merchant_entity.SearchMerchantResponse{
		Data: merchantSearched,
		Meta: &merchant_entity.MetaData{
			Limit:  searchQuery.Limit,
			Offset: searchQuery.Offset,
			Total:  len(*merchantSearched),
		},
	}, nil

}
