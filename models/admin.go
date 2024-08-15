package models

type Admin struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Phone_number string `json:"phone_number,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	Address      string `json:"addres,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	DeletedAt    string `json:"delete_at,omitempty"`
}

type AdminCreate struct {
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"addres"`
}

type AdminUpdate struct {
	Id           string `json:"-"`
	Name         string `json:"name"`
	Phone_number string `json:"phone_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"addres"`
}

type AdminPrimaryKey struct {
	Id string `json:"id"`
}

type AdminGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type AdminGetListResponse struct {
	Count int      `json:"count"`
	Admin []*Admin `json:"admin"`
}
