package rest

import (
	"encoding/json"
	"net/http"
	"ordersystem/model"
	"ordersystem/repository"

	"github.com/go-chi/render"
)

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		// get slice from db
		drinks := db.GetDrinks()

		render.Status(r, http.StatusOK)
		render.JSON(w, r, drinks)
	}
}

// GetOrders      godoc
// @tags          Order
// @Description   Returns all orders
// @Produce       json
// @Success       200 {array} model.Order
// @Router        /api/order/all [get]
// todo create GetOrders /api/order/all
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders := db.GetOrders()

		render.Status(r, http.StatusOK)
		render.JSON(w, r, orders)
	}
}

// GetOrdersTotal godoc
// @tags          Order
// @Description   Returns total counts for ordered drinks
// @Produce       json
// @Success       200 {object} map[string]int
// @Router        /api/order/totalled [get]
// todo create GetOrdersTotal /api/order/total
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totals := db.GetTotalledOrders()

		render.Status(r, http.StatusOK)
		render.JSON(w, r, totals)
	}
}

// PostOrder 		godoc
// @tags 			Order
// @Description 	Adds an order to the db
// @Accept 			json
// @Param 			b body model.Order true "Order"
// @Produce  		json
// @Success 		200
// @Failure     	400
// @Router 			/api/order [post]
func PostOrder(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		// declare empty order struct
		var newOrder model.Order

		err := json.NewDecoder(r.Body).Decode(&newOrder)

		// handle error and render Status 400
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Invalid request")
		}
		// add to db
		db.AddOrder(&newOrder)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
	}
}
