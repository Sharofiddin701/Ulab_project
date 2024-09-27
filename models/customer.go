package models

type Customer struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Phone_number string `json:"phone_number"`
	Birthday     string `json:"birthday"`
	Gender       string `json:"gender"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	DeletedAt    string `json:"delete_at,omitempty"`
}

type CustomerCreate struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Phone_number string `json:"phone_number"`
	Birthday     string `json:"birthday"`
	Gender       string `json:"gender"`
}

type CustomerUpdate struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Phone_number string `json:"phone_number"`
	Birthday     string `json:"birthday"`
	Gender       string `json:"gender"`
}

type CustomerPrimaryKey struct {
	Id string `json:"id"`
}

type CustomerGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CustomerGetListResponse struct {
	Count    int         `json:"count"`
	Customer []*Customer `json:"customer"`
}
