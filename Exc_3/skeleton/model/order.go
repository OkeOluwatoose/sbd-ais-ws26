package model

import "time"

type Order struct {
	DrinkID uint64 `json:"drink_id"` // foreign key
	// todo Add fields: CreatedAt (time.Time), Amount with suitable types
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`

	// todo json attributes need to be snakecase, i.e. name, created_at, my_variable, ..
}
