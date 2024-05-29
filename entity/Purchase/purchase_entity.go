package purchase_entity

type Purchase struct {
	Id                             string  `json:"purchaseId,omitempty"`
	UserId                         string  `json:"userId,omitempty"`
	TotalPrice                     int     `json:"totalPrice,omitempty"`
	EstimatedDeliveryTimeInMinutes int     `json:"estimatedDeliveryTimeInMinutes,omitempty"`
	CalculatedEstimateId           string  `json:"calculatedEstimateId,omitempty"`
	Status                         string  `json:"status,omitempty"`
	Distance                       float32 `json:"distance,omitempty"`
	Latitude                       float64 `json:"lat,omitempty"`
	Longitude                      float64 `json:"long,omitempty"`
	Order                          []Order `json:"orders,omitempty"`
	CreatedAt                      string  `json:"createdAt,omitempty"`
}

type Order struct {
	Id              string      `json:"orderId,omitempty"`
	MechantId       string      `json:"merchantId,omitempty" validate:"required"`
	PurchaseId      string      `json:"purchaseId,omitempty"`
	IsStartingPoint *bool       `json:"isStartingPoint,omitempty" validate:"required"`
	OrderItem       []OrderItem `json:"items,omitempty" validate:"required,dive"`
	CreatedAt       string      `json:"createdAt,omitempty"`
}

type OrderItem struct {
	Id              string `json:"Id,omitempty"`
	PurchaseOrderId string `json:"purchaseOrderId,omitempty"`
	ItemId          string `json:"itemId,omitempty" validate:"required"`
	Quantity        int    `json:"quantity,omitempty" validate:"required"`
	CreatedAt       string `json:"createdAt,omitempty"`
}
