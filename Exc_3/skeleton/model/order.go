package model

import "time"

type Order struct {
	DrinkID uint64 `json:"drink_id"` // foreign key
	// todo Add fields: CreatedAt (time.Time), Amount with suitable types
	amount     int
	created_at time.Time

	// todo json attributes need to be snakecase, i.e. name, created_at, my_variable, ..
}
