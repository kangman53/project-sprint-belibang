package merchant_entity

type AddMerchantResponse struct {
	Id string `json:"merchantId"`
}

type SearchMerchantResponse struct {
	Data *[]SearchMerchantData `json:"data"`
	Meta *MetaData             `json:"meta"`
}

type SearchMerchantData struct {
	Id        string        `json:"merchantId"`
	Name      string        `json:"name"`
	Category  string        `json:"merchantCategory"`
	ImageUrl  string        `json:"imageUrl"`
	Location  *LocationData `json:"location"`
	CreatedAt string        `json:"createdAt"`
}

type LocationData struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type MetaData struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
