package models

type Product struct {
	Id                  string `json:"id,omitempty"`
	Is_favourite        bool   `json:"is_favourite,omitempty"`
	Image               string `json:"image,omitempty"`
	Name                string `json:"name,omitempty"`
	Product_categoty    string `json:"product_categoty,omitempty"`
	Price               int    `json:"price,omitempty"`
	Price_with_discount int    `json:"price_with_discount,omitempty"`
	Rating              int    `json:"rating,omitempty"`
	Order_count         int    `json:"order_count,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
	DeletedAt           string `json:"delete_at,omitempty"`
}

type ProductCreate struct {
	Is_favourite        bool   `json:"is_favourite"`
	Image               string `json:"image"`
	Name                string `json:"name"`
	Product_categoty    string `json:"product_categoty"`
	Price               int    `json:"price"`
	Price_with_discount int    `json:"price_with_discount"`
	Rating              int    `json:"rating"`
	Order_count         int    `json:"order_count"`
}

type ProductUpdate struct {
	Id                  string `json:"id"`
	Is_favourite        bool   `json:"is_favourite"`
	Image               string `json:"image"`
	Name                string `json:"name"`
	Product_categoty    string `json:"product_categoty"`
	Price               int    `json:"price"`
	Price_with_discount int    `json:"price_with_discount"`
	Rating              int    `json:"rating"`
	Order_count         int    `json:"order_count"`
}

type ProductPrimaryKey struct {
	Id string `json:"id"`
}
type ProductGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ProductGetListResponse struct {
	Count   int       `json:"count"`
	Product []Product `json:"product"`
}
