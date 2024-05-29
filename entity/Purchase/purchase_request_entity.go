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
