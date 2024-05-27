package item_entity

type AddItemRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Category string `json:"productCategory" validate:"required,validateCategory=item"`
	Price    int    `json:"price" validate:"required,min=1"`
	ImageUrl string `json:"imageUrl" validate:"required,validateUrl"`
}

type SearchItemQuery struct {
	ItemId, Name, Category, CreatedAt string
	Limit, Offset                     int
}
