package domain

type Car struct {
	Model        string `json:"model"`
	Registration string `json:"registration"`
	Mileage      int    `json:"mileage"`
	IsRented     bool   `json:"isrented"`
}
