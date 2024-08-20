package models

type Favorite struct {
	ProductID string `json:"product_id"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"delete_at,omitempty"`
}

type FavoriteGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type FavoriteGetListResponse struct {
	Total    int        `json:"total"`
	Favorite []Favorite `json:"favorite"`
}
