package repository

import "ordersystem/model"

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

// todo
func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	drinks := []model.Drink{
		{ID: 1, Description: "Rich espresso combined with milk, vanilla-flavored syrup, and a drizzle of sweet caramel sauce served over ice.", Name: "Iced Caramel Macchiato", Price: 5.45},
		{ID: 2, Description: "Smooth and creamy pure Japanese matcha green tea whisked with steamed milk and lightly sweetened.", Name: "Classic Matcha Latte", Price: 4.95},
		{ID: 3, Description: "Slow-steeped specialty coffee blend brewed in cool water for 20 hours to create an ultra-smooth, bold flavor.", Name: "Cold Brew Coffee", Price: 4.25},
		{ID: 4, Description: "A sweet, vibrant beverage featuring strawberry flavors, accents of passion fruit, and real fruit pieces shaken with ice.", Name: "Strawberry Acai Refresher", Price: 4.75},
	}

	// Init orders slice with some test data
	orders := []model.Order{
		{DrinkID: 1},
		{DrinkID: 4},
		{DrinkID: 2},
	}
	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

// todo
func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	totalledOrders := make(map[uint64]uint64)
	// key = DrinkID, value = Amount of orders
	for _, order := range db.orders {
		totalledOrders[order.DrinkID]++
	}
	// totalledOrders map[uint64]uint64
	return totalledOrders
}

func (db *DatabaseHandler) AddOrder(order *model.Order) {
	// todo
	// add order to db.orders slice
	db.orders = append(db.orders, *order)
}
