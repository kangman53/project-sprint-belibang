package purchase_entity

type PurchaseEstimateResponse struct {
	TotalPrice                     int     `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int     `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId           string  `json:"calculatedEstimateId"`
	Distance                       float32 `json:"distance"`
}
