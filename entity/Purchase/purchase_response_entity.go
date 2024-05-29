package purchase_entity

import (
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type PurchaseEstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string `json:"calculatedEstimateId"`
}

type PurchaseOrderResponse struct {
	OrderId string `json:"orderId"`
}

type SearchHistoryOrderResponse struct {
	OrderId string                      `json:"orderId"`
	Orders  *[]SearchHistoryOrderOrders `json:"orders"`
}

type SearchHistoryOrderOrders struct {
	Merchant *merchant_entity.MerchantData `json:"merchant"`
	Items    *[]item_entity.SearchItemData `json:"items"`
}
