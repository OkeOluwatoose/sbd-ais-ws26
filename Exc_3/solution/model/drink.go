package model

type Drink struct {
	ID uint64 `json:"id"`
	// todo Add fields: Name, Price, Description with suitable types
	Description string
	Name        string
	Price       float64
	// todo json attributes need to be snakecase, i.e. name, created_at, my_variable, ..

}
