package models

type Color struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Name      string `json:"color_name"`
	Url       string `json:"color_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"delete_at"`
}

type ColorCreate struct {
	ProductId string `json:"product_id"`
	Name      string `json:"color_name"`
	Url       string `json:"color_url"`
}

type ColorUpdate struct {
	ProductId string `json:"product_id"`
	Id        string `json:"id"`
	Name      string `json:"color_name"`
	Url       string `json:"color_url"`
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
