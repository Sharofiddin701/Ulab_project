package models

type Order struct {
	Id         string       `json:"id,omitempty"`
	CustomerId string       `json:"customer_id,omitempty"`
	TotalPrice float64      `json:"total_price,omitempty"`
	Status     string       `json:"status,omitempty"`
	CreatedAt  string       `json:"created_at,omitempty"`
	UpdatedAt  string       `json:"updated_at,omitempty"`
	DeletedAt  string       `json:"delete_at,omitempty"`
	OrderItems []OrderItems `json:"order_items,omitempty"`
}

type OrderCreate struct {
	CustomerId string `json:"customer_id"`
}

type OrderUpdate struct {
	Id         string  `json:"id"`
	CustomerId string  `json:"customer_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type OrderGetListRequest struct {
	CustomerId string `json:"customer_id"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type OrderGetListResponse struct {
	Count int      `json:"count"`
	Order []*Order `json:"order"`
}

type OrderCreateRequest struct {
	Order Order  `json:"order"`
	Items []OrderItems `json:"items"`
}

type SwaggerOrderCreateRequest struct {
	Order OrderCreate         `json:"order"`
	Items []SwaggerOrderItems `json:"items"`
}
