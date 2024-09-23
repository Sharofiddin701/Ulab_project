package models

type Color struct {
	Id        string   `json:"id"`
	ProductId string   `json:"product_id"`
	Name      string   `json:"color_name"`
	Url       []string `json:"color_url"`
	Count     int      `json:"count"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	DeletedAt string   `json:"delete_at,omitempty"`
}

type ColorCreate struct {
	ProductId string   `json:"product_id"`
	Name      string   `json:"color_name"`
	Url       []string `json:"color_url"`
	Count     int      `json:"count"`
}

type ColorUpdate struct {
	ProductId string   `json:"product_id"`
	Id        string   `json:"id"`
	Name      string   `json:"color_name"`
	Url       []string `json:"color_url"`
	Count     int      `json:"count"`
}
type ColorPrimaryKey struct {
	Id string `json:"id"`
}

type ColorGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ColorGetListResponse struct {
	Count int      `json:"count"`
	Color []*Color `json:"color"`
}
