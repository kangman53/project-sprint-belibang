package merchant_entity

type MerchantRegisterRequest struct {
	Name     string           `json:"name" validate:"required,min=2,max=30"`
	Category string           `json:"merchantCategory" validate:"required,merchantCategory"`
	ImageUrl string           `json:"imageUrl" validate:"required,validateUrl"`
	Location *LocationDetails `json:"location" validate:"required"`
}

type LocationDetails struct {
	Latitude  float64 `json:"latitude" validate:"required,latitudeValidation"`
	Longitude float64 `json:"longitude" validate:"required,longitudeValidation"`
}
