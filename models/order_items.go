package models

type OrderItems struct {
	Id         string  `json:"id,omitempty"`
	OrderId    string  `json:"order_id,omitempty"`
	ProductId  string  `json:"product_id,omitempty"`
	Quantity   int     `json:"quantity,omitempty"`
	Price      float64 `json:"price,omitempty"`
	TotalPrice float64 `json:"total,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
	DeletedAt  string  `json:"delete_at,omitempty"`
}

type SwaggerOrderItems struct {
	ProductId string `json:"product_id,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
}

type OrderItemsCreate struct {
	OrderId    string  `json:"order_id"`
	ProductId  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}

type OrderItemsUpdate struct {
	Id         string  `json:"id"`
	OrderId    string  `json:"order_id"`
	ProductId  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}

type OrderItemsPrimaryKey struct {
	Id string `json:"id"`
}

type OrderItemsGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type OrderItemsGetListResponse struct {
	Count int           `json:"count"`
	Items []*OrderItems `json:"order_items"`
}
