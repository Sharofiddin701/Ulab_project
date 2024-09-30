package models

type Order struct {
	Id             string       `json:"id,omitempty"`
	CustomerId     string       `json:"customer_id,omitempty"`
	AddressName    string       `json:"address_name,omitempty"`
	Longtitude     float64      `json:"longtitude"`
	Latitude       float64      `json:"latitude"`
	TotalPrice     float64      `json:"total_price,omitempty"`
	Status         string       `json:"status,omitempty"`
	DeliveryStatus string       `json:"delivery_status,omitempty"`
	DeliveryCost   float64      `json:"delivery_cost,omitempty"`
	PaymentMethod  string       `json:"payment_method,omitempty"`
	PaymentStatus  string       `json:"payment_status,omitempty"`
	CreatedAt      string       `json:"created_at,omitempty"`
	UpdatedAt      string       `json:"updated_at,omitempty"`
	DeletedAt      string       `json:"delete_at,omitempty"`
	OrderItems     []OrderItems `json:"order_items,omitempty"`
}

type OrderCreate struct {
	CustomerId     string  `json:"customer_id"`
	AddressName    string  `json:"address_name,omitempty"`
	Longtitude     float64 `json:"longtitude"`
	Latitude       float64 `json:"latitude"`
	DeliveryStatus string  `json:"delivery_status"`
	DeliveryCost   float64 `json:"delivery_cost"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentStatus  string  `json:"payment_status"`
}

type OrderUpdate struct {
	Id             string  `json:"id"`
	CustomerId     string  `json:"customer_id"`
	AddressName    string  `json:"address_name,omitempty"`
	Longtitude     float64 `json:"longtitude"`
	Latitude       float64 `json:"latitude"`
	DeliveryStatus string  `json:"delivery_status"`
	DeliveryCost   float64 `json:"delivery_cost"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentStatus  string  `json:"payment_status"`
	TotalPrice     float64 `json:"total_price"`
	Status         string  `json:"status"`
}

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type OrderGetListRequest struct {
	CustomerId     string  `json:"customer_id"`
	Longtitude     float64 `json:"longtitude"`
	Latitude       float64 `json:"latitude"`
	AddressName    string  `json:"address_name,omitempty"`
	DeliveryStatus string  `json:"delivery_status"`
	DeliveryCost   float64 `json:"delivery_cost"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentStatus  string  `json:"payment_status"`
	Offset         int     `json:"offset"`
	Limit          int     `json:"limit"`
}

type OrderGetListResponse struct {
	Count int      `json:"count"`
	Order []*Order `json:"order"`
}

type OrderCreateRequest struct {
	Order Order        `json:"order"`
	Items []OrderItems `json:"items"`
}

type SwaggerOrderCreateRequest struct {
	Order OrderCreate         `json:"order"`
	Items []SwaggerOrderItems `json:"items"`
}
