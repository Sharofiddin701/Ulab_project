package models

type OrderProduct struct {
	Id        string `json:"id,omitempty"`
	OrderId   string `json:"order_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
	Price     string `json:"price,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"delete_at,omitempty"`
}

type OrderProductCreate struct {
	OrderId   string `json:"order_id"`
	ProductId string `json:"product_id"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
}

type OrderProductUpdate struct {
	Id        string `json:"id"`
	OrderId   string `json:"order_id"`
	ProductId string `json:"product_id"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
}

type OrderProductPrimaryKey struct {
	Id string `json:"id"`
}

type OrderProductGetListRequest struct {
	OrderId   string `json:"order_id"`
	ProductId string `json:"product_id"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type OrderProductGetListResponse struct {
	Count        int             `json:"count"`
	OrderProduct []*OrderProduct `json:"order_product"`
}
