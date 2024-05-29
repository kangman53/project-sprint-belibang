package merchant_service

import (
	"fmt"
	"strconv"
	"strings"

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

func (service *merchantServiceImpl) SearchNearby(ctx *fiber.Ctx, query merchant_entity.SearchNearbyMerchantQuery) (merchant_entity.SearchNearbyMerchantResponse, error) {
	if err := service.Validator.Struct(query); err != nil {
		return merchant_entity.SearchNearbyMerchantResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	params := strings.Split(ctx.Params("coordinate"), ",")
	if len(params) != 2 {
		return merchant_entity.SearchNearbyMerchantResponse{}, exc.BadRequestException("latitude or longitude not valid")
	}

	latParam := params[0]
	longParam := params[1]

	// latParam := ctx.Params("lat")
	// longParam := ctx.Params("long")

	lat, err1 := strconv.ParseFloat(latParam, 64)
	long, err2 := strconv.ParseFloat(longParam, 64)

	if err1 != nil || err2 != nil {
		return merchant_entity.SearchNearbyMerchantResponse{}, exc.BadRequestException("latitude or longitude not valid")
	}

	userCtx := ctx.UserContext()
	dataSearch, err := service.MerchantRepository.SearchNearby(userCtx, query, lat, long)
	if err != nil {
		return merchant_entity.SearchNearbyMerchantResponse{}, err
	}

	return merchant_entity.SearchNearbyMerchantResponse{
		Data: dataSearch,
	}, nil
}
