package item_entity

type AddItemResponse struct {
	Id string `json:"itemId"`
}

type SearchItemResponse struct {
	Data []*SearchItemData `json:"data"`
	Meta *MetaData         `json:"meta"`
}

type SearchItemData struct {
	Id        string `json:"itemId"`
	Name      string `json:"name"`
	Category  string `json:"productCategory"`
	Price     int    `json:"price"`
	ImageUrl  string `json:"imageUrl"`
	CreatedAt string `json:"createdAt"`
}

type MetaData struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
