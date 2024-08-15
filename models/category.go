package models

type Category struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"delete_at,omitempty"`
}

type CategoryCreate struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CategoryUpdate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CategoryPrimaryKey struct {
	Id string `json:"id"`
}

type CategoryGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CategoryGetListResponse struct {
	Count    int         `json:"count"`
	Category []*Category `json:"category"`
}
