package models

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentId  string `json:"parent_id"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"delete_at,omitempty"`
}

type CategoryCreate struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	ParentId string `json:"parent_id"`
}

type CategoryUpdate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
