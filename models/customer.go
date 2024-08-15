package models

type Customer struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Phone_number string `json:"phone_number,omitempty"`
	Address      string `json:"address,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	DeletedAt    string `json:"delete_at,omitempty"`
}

type CustomerCreate struct {
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type CustomerUpdate struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Password     string `json:"password"`
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
