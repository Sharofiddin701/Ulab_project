package models

type Banner struct {
	Id           string `json:"id,omitempty"`
	Banner_image string `json:"banner_image,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	DeletedAt    string `json:"delete_at,omitempty"`
}

type BannerCreate struct {
	Banner_image string `json:"banner_image"`
}

type BannerUpdate struct {
	Id           string `json:"id"`
	Banner_image string `json:"banner_image"`
}

type BannerPrimaryKey struct {
	Id string `json:"id"`
}

type BannerGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type BannerGetListResponse struct {
	Count  int       `json:"count"`
	Banner []*Banner `json:"banner"`
}
