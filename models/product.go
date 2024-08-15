package models

type Product struct {
	Id             string `json:"id,omitempty"`
	CategoryId     string `json:"category_id,omitempty"`
	BrandId        string `json:"brand_id,omitempty"`
	Name           string `json:"name,omitempty"`
	ProductArticle string `json:"product_articl"`
	Count          int    `json:"count,omitempty"`
	Price          int    `json:"price,omitempty"`
	ProductImage   string `json:"product_image,omitempty"`
	Comment        string `json:"comment,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	DeletedAt      string `json:"delete_at,omitempty"`
}

type ProductCreate struct {
	CategoryId     string `json:"category_id"`
	BrandId        string `json:"brand_id"`
	Name           string `json:"name"`
	ProductArticle string `json:"product_articl"`
	Count          int    `json:"count"`
	Price          int    `json:"price"`
	ProductImage   string `json:"product_image"`
	Comment        string `json:"comment"`
}

type ProductUpdate struct {
	Id             string `json:"id"`
	CategoryId     string `json:"category_id"`
	BrandId        string `json:"brand_id"`
	Name           string `json:"name"`
	ProductArticle string `json:"product_articl"`
	Count          int    `json:"count"`
	Price          int    `json:"price"`
	ProductImage   string `json:"product_image"`
	Comment        string `json:"comment"`
}

type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type ProductGetListRequest struct {
	CategoryId string `json:"category_id"`
	BrandId    string `json:"brand_id"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type ProductGetListResponse struct {
	Count   int        `json:"count"`
	Product []*Product `json:"product"`
}
