package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,number,min=0"`
	Description string  `json:"description"`
}
