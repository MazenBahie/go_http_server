package models

type CreditCard struct {
	ID             int    `json:"id"`
	Name           string `json:"name" validate:"required"`
	CardNumber     string `json:"card_number" validate:"required,number"`
	Cvv            string `json:"cvv" validate:"required,number"`
	ExpirationDate string `json:"expiration_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	UserID         int    `json:"user_id" validate:"required,number,gt=0"`
}
