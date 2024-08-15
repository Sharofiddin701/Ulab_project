package models

type Brand struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Brand_image string `json:"brand_image,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"delete_at,omitempty"`
}

type BrandCreate struct {
	Name        string `json:"name"`
	Brand_image string `json:"brand_image"`
}

type BrandUpdate struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Brand_image string `json:"brand_image"`
}

type BrandPrimaryKey struct {
	Id string `json:"id"`
}

type BrandGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type BrandGetListResponse struct {
	Count int      `json:"count"`
	Brand []*Brand `json:"brand"`
}
