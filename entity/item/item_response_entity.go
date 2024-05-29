package item_entity

type AddItemResponse struct {
	Id string `json:"itemId"`
}

type SearchItemData struct {
	Id        string `json:"itemId"`
	Name      string `json:"name"`
	Category  string `json:"productCategory"`
	Price     int    `json:"price"`
	ImageUrl  string `json:"imageUrl"`
	CreatedAt string `json:"createdAt"`
}
