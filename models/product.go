package models

type Product struct {
	Id              string  `json:"id,omitempty"`
	CategoryId      string  `json:"category_id,omitempty"`
	Favorite        bool    `json:"favorite"`
	Name            string  `json:"name,omitempty"`
	Price           float64 `json:"price,omitempty"`
	WithDiscount    float64 `json:"with_discount"`
	Rating          float64 `json:"rating,omitempty"`
	Description     string  `json:"description,omitempty"`
	OrderCount      int     `json:"order_count,omitempty"`
	Color           []Color `json:"color,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
	DeletedAt       string  `json:"deleted_at,omitempty"`
	Status          string  `json:"status"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountEndTime string  `json:"discount_end_time"`
}

type ProductCreate struct {
	CategoryId      string  `json:"category_id"`
	Favorite        bool    `json:"favorite"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	With_discount   float64 `json:"with_discount"`
	Rating          float64 `json:"rating"`
	Description     string  `json:"description"`
	Order_count     int     `json:"order_count"`
	Status          string  `json:"status"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountEndTime string  `json:"discount_end_time"`
}

type ProductUpdate struct {
	Id              string  `json:"id"`
	CategoryId      string  `json:"category_id"`
	Favorite        bool    `json:"favorite"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	With_discount   float64 `json:"with_discount"`
	Rating          float64 `json:"rating"`
	Description     string  `json:"description"`
	Order_count     int     `json:"order_count"`
	Status          string  `json:"status"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountEndTime string  `json:"discount_end_time"`
}

type ProductPrimaryKey struct {
	Id string `json:"id"`
}
type ProductGetListRequest struct {
	CategoryId string `json:"category_id"`
	Favorite   *bool  `json:"favorite"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Name       string `json:"name"`
}

type ProductGetListResponse struct {
	Count   int       `json:"count"`
	Product []Product `json:"product"`
}
