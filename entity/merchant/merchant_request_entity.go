package merchant_entity

type AddMerchantRequest struct {
	Name     string           `json:"name" validate:"required,min=2,max=30"`
	Category string           `json:"merchantCategory" validate:"required,validateCategory=merchant"`
	ImageUrl string           `json:"imageUrl" validate:"required,validateUrl"`
	Location *LocationDetails `json:"location" validate:"required"`
}

type LocationDetails struct {
	Latitude  float64 `json:"latitude" validate:"required,validateGeoCoord=lat"`
	Longitude float64 `json:"longitude" validate:"required,validateGeoCoord=long"`
}
