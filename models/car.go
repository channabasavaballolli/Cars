package models

type Car struct {
	ID      int     `json:"id"`
	Make    string  `json:"make"`
	Model   string  `json:"model"`
	Year    int     `json:"year"`
	Price   float64 `json:"price"`
	Color   string  `json:"color"`
	Mileage int     `json:"mileage"`
}
