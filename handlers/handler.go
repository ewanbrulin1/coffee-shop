package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ewanb/CoffeeGo/coffee-shop-api/model" // Importing model for data types
	"github.com/gorilla/mux"
)

var orderCounter int = 1

// GetMenu : Handler exporté (GET /menu)
func GetMenu(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, model.Drinks)
}

// GetDrink : Handler exporté (GET /menu/{id})
func GetDrink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	drinkID := vars["id"]

	for _, drink := range model.Drinks {
		if drink.ID == drinkID {
			RespondJSON(w, http.StatusOK, drink)
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Boisson non trouvée")
}

// FindDrinkByID : Helper exporté (utilisé par order.go)
func FindDrinkByID(id string) (model.Drink, bool) {
	for _, drink := range model.Drinks {
		if drink.ID == id {
			return drink, true
		}
	}
	return model.Drink{}, false
}

// CreateOrder : Handler exporté
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		RespondError(w, http.StatusBadRequest, "Corps de requête invalide")
		return
	}

	foundDrink, drinkFound := FindDrinkByID(newOrder.DrinkID)
	if !drinkFound {
		RespondError(w, http.StatusBadRequest, "ID de boisson introuvable")
		return
	}

	newOrder.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	newOrder.DrinkName = foundDrink.Name
	newOrder.Status = model.StatusPending
	newOrder.OrderedAt = time.Now()
	newOrder.TotalPrice = model.CalculatePrice(foundDrink.BasePrice, newOrder.Size, newOrder.Extras)

	model.Orders = append(model.Orders, newOrder)

	RespondJSON(w, http.StatusCreated, newOrder)
}

// GetOrders : Handler exporté
func GetOrders(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, model.Orders)
}

// GetOrder : Handler exporté
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for _, order := range model.Orders {
		if order.ID == orderID {
			RespondJSON(w, http.StatusOK, order)
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée")
}

// UpdateOrderStatus : Handler exporté
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var statusUpdate struct {
		Status model.OrderStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		RespondError(w, http.StatusBadRequest, "Corps de requête invalide pour le statut")
		return
	}

	for i := range model.Orders {
		if model.Orders[i].ID == orderID {
			model.Orders[i].Status = statusUpdate.Status
			RespondJSON(w, http.StatusOK, model.Orders[i])
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée")
}

// DeleteOrder : Handler exporté
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for i, order := range model.Orders {
		if order.ID == orderID {
			if order.Status == model.StatusPickedUp {
				RespondError(w, http.StatusBadRequest, "Impossible d'annuler une commande déjà récupérée")
				return
			}
			model.Orders = append(model.Orders[:i], model.Orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée")
}
