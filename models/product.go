package models

type Product struct {
	Id                  string  `json:"id"`
	Favorite            bool    `json:"favorite"`
	Image               string  `json:"image"`
	Name                string  `json:"name"`
	Product_categoty    string  `json:"product_categoty"`
	Price               int     `json:"price"`
	Price_with_discount int     `json:"price_with_discount"`
	Rating              float64 `json:"rating"`
	Description         string  `json:"description"`
	Order_count         int     `json:"order_count"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	DeletedAt           string  `json:"delete_at"`
}

type ProductCreate struct {
	Favorite            bool    `json:"favorite"`
	Image               string  `json:"image"`
	Name                string  `json:"name"`
	Product_categoty    string  `json:"product_categoty"`
	Price               int     `json:"price"`
	Price_with_discount int     `json:"price_with_discount"`
	Rating              float64 `json:"rating"`
	Description         string  `json:"description"`
	Order_count         int     `json:"order_count"`
}

type ProductUpdate struct {
	Id                  string  `json:"id"`
	Favorite            bool    `json:"favorite"`
	Image               string  `json:"image"`
	Name                string  `json:"name"`
	Product_categoty    string  `json:"product_categoty"`
	Price               int     `json:"price"`
	Price_with_discount int     `json:"price_with_discount"`
	Rating              float64 `json:"rating"`
	Description         string  `json:"description"`
	Order_count         int     `json:"order_count"`
}

type ProductPrimaryKey struct {
	Id string `json:"id"`
}
type ProductGetListRequest struct {
	Favorite *bool `json:"favorite"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type ProductGetListResponse struct {
	Count   int       `json:"count"`
	Product []Product `json:"product"`
}
