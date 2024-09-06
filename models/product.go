package models

type Product struct {
	Id            string   `json:"id,omitempty"`
	CategoryId    string   `json:"category_id,omitempty"`
	Favorite      bool     `json:"favorite"`
	Image         []string `json:"image,omitempty"`
	Name          string   `json:"name,omitempty"`
	Price         int      `json:"price,omitempty"`
	With_discount int      `json:"with_discount,omitempty"`
	Rating        float64  `json:"rating,omitempty"`
	Description   string   `json:"description,omitempty"`
	Order_count   int      `json:"order_count,omitempty"`
	Color         []Color  `json:"color,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
	DeletedAt     string   `json:"delete_at,omitempty"`
}

type ProductCreate struct {
	CategoryId    string   `json:"category_id"`
	Favorite      bool     `json:"favorite"`
	Image         []string `json:"image"`
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	With_discount int      `json:"with_discount"`
	Rating        float64  `json:"rating"`
	Description   string   `json:"description"`
	Order_count   int      `json:"order_count"`
}

type ProductUpdate struct {
	Id            string   `json:"id"`
	CategoryId    string   `json:"category_id"`
	Favorite      bool     `json:"favorite"`
	Image         []string `json:"image"`
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	With_discount int      `json:"with_discount"`
	Rating        float64  `json:"rating"`
	Description   string   `json:"description"`
	Order_count   int      `json:"order_count"`
}

type ProductPrimaryKey struct {
	Id string `json:"id"`
}
type ProductGetListRequest struct {
	CategoryId string `json:"category_id"`
	Favorite   *bool  `json:"favorite"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type ProductGetListResponse struct {
	Count   int       `json:"count"`
	Product []Product `json:"product"`
}
