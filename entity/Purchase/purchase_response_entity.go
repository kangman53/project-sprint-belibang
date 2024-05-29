package purchase_entity

type PurchaseEstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string `json:"calculatedEstimateId"`
}

type PurchaseOrderResponse struct {
	OrderId string `json:"orderId"`
}
