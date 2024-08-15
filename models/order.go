package models

type Order struct {
	Id         string `json:"id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	Shipping   string `json:"shipping,omitempty"`
	Payment    string `json:"payment,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	DeletedAt  string `json:"delete_at,omitempty"`
}

type OrderCreate struct {
	CustomerId string `json:"customer_id"`
	Shipping   string `json:"shipping"`
	Payment    string `json:"payment"`
}

type OrderUpdate struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	Shipping   string `json:"shipping"`
	Payment    string `json:"payment"`
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
	Order []*Order `json:"brand"`
}
