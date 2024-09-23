package models

type Location struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Info      string  `json:"info"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Image     string  `json:"image"`
}


type LocationCreate struct {
	Name      string  `json:"name"`	
	Info      string  `json:"info"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Image     string  `json:"image"`
}

type LocationUpdate struct {
	
}