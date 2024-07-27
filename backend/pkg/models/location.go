package models

// Location models a location.
type Location struct {
	Latitude   float32 `json:"latitude"`
	Longitutde float32 `json:"longitutde"`
	ID         int     `json:"id"`
	Location   string  `json:"location"`
	Image      string  `json:"image"`
}
