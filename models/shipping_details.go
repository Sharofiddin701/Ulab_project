package models

type ShippingDetails struct {
	Id              string  `json:"id"`
	OrderId         string  `json:"order_id"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	PostalCode      string  `json:"postal_code"`
	Phone_number    string  `json:"phone_number"`
	Delivery_method string  `json:"deleviry_status"`
	Delivery_cost   float64 `json:"delivery_cost"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
	Deleted_at      string  `json:"deleted_at"`
}

type ShippingDetailsCreate struct {
	OrderId         string  `json:"order_id"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	PostalCode      string  `json:"postal_code"`
	Phone_number    string  `json:"phone_number"`
	Delivery_method string  `json:"deleviry_status"`
	Delivery_cost   float64 `json:"delivery_cost"`
}

type ShippingDetailsUpdate struct {
	Id              string  `json:"id"`
	OrderId         string  `json:"order_id"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	PostalCode      string  `json:"postal_code"`
	Phone_number    string  `json:"phone_number"`
	Delivery_method string  `json:"deleviry_status"`
	Delivery_cost   float64 `json:"delivery_cost"`
}

type ShippingDetailsPrimaryKey struct {
	Id string `json:"id"`
}

type ShippingDetailsGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ShippingDetailsGetListResponse struct {
	Count int                `json:"count"`
	Items []*ShippingDetails `json:"shipping_details"`
}
