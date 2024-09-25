package models

type Location struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Info      string  `json:"info"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Image     string  `json:"image"`
	OpensAt   string  `json:"opens_at"`
	ClosesAt  string  `json:"closes_at"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
	DeletedAt string  `json:"delete_at,omitempty"`
}

type LocationCreate struct {
	Name      string  `json:"name"`
	Info      string  `json:"info"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Image     string  `json:"image"`
	OpensAt   string  `json:"opens_at"`
	ClosesAt  string  `json:"closes_at"`
}

type LocationUpdate struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Info      string  `json:"info"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Image     string  `json:"image"`
	OpensAt   string  `json:"opens_at"`
	ClosesAt  string  `json:"closes_at"`
}

type LacationPrimaryKey struct {
	Id string `json:"id"`
}

type LocationGetListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type LocationGetListResponse struct {
	Count    int         `json:"count"`
	Location []*Location `json:"locations"`
}
