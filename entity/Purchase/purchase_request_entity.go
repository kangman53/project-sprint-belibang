package purchase_entity

import (
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type PurchaseEstimateRequest struct {
	UserLocation *merchant_entity.LocationDetails `json:"userLocation" validate:"required"`
	Orders       *[]Order                         `json:"orders" validate:"required,dive"`
}

type PurchaseOrderRequest struct {
	CalculatedEstimateId string `json:"calculatedEstimateId" validate:"required"`
}

type SearcHistoryOrderQuery struct {
	MerchantId       string `query:"merchantId"`
	Name             string `query:"name"`
	MerchantCategory string `query:"merchantCategory"`
	Limit            int    `query:"limit"`
	Offset           int    `query:"offset"`
}
