package merchant_entity

type AddMerchantResponse struct {
	Id string `json:"merchantId"`
}

type SearchNearbyMerchantResponse struct {
	Data *[]SearchNearbyMerchantData `json:"data"`
}
type SearchNearbyMerchantData struct {
	Merchant *MerchantData `json:"merchant"`
	Items    *[]Item       `json:"items"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantData struct {
	Id               string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageURL         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
}

type Item struct {
	ItemID          string  `json:"itemId"`
	Name            string  `json:"name"`
	ProductCategory string  `json:"productCategory"`
	Price           float64 `json:"price"`
	ImageURL        string  `json:"imageUrl"`
	CreatedAt       string  `json:"createdAt"`
}
