package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ⚠️ Plus d'importation problématique du module racine.

// OrderStatus : structure exportée
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked-up"
	StatusCancelled OrderStatus = "cancelled"
)

// Order : structure exportée
type Order struct {
	ID           string      `json:"id"`
	DrinkID      string      `json:"drink_id"`
	DrinkName    string      `json:"drink_name"`
	Size         string      `json:"size"`
	Extras       []string    `json:"extras"`
	CustomerName string      `json:"customer_name"`
	Status       OrderStatus `json:"status"`
	TotalPrice   float64     `json:"total_price"`
	OrderedAt    time.Time   `json:"ordered_at"`
}

// Orders : données exportées
var Orders []Order
var orderCounter int = 1

// CalculatePrice : fonction exportée
func CalculatePrice(basePrice float64, size string, extras []string) float64 {
	price := basePrice
	switch size {
	case "small":
		price *= 0.8
	case "medium":
	case "large":
		price *= 1.3
	}
	price += float64(len(extras)) * 0.50
	return price
}

// CreateOrder : Handler exporté
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		RespondError(w, http.StatusBadRequest, "Corps de requête invalide") // ✅ Correction
		return
	}

	foundDrink, drinkFound := FindDrinkByID(newOrder.DrinkID)
	if !drinkFound {
		RespondError(w, http.StatusBadRequest, "ID de boisson introuvable") // ✅ Correction
		return
	}

	newOrder.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	newOrder.DrinkName = foundDrink.Name
	newOrder.Status = StatusPending
	newOrder.OrderedAt = time.Now()
	newOrder.TotalPrice = CalculatePrice(foundDrink.BasePrice, newOrder.Size, newOrder.Extras)

	Orders = append(Orders, newOrder)

	RespondJSON(w, http.StatusCreated, newOrder) // ✅ Correction
}

// GetOrders : Handler exporté
func GetOrders(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, Orders) // ✅ Correction
}

// GetOrder : Handler exporté
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for _, order := range Orders {
		if order.ID == orderID {
			RespondJSON(w, http.StatusOK, order) // ✅ Correction
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée") // ✅ Correction
}

// UpdateOrderStatus : Handler exporté
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var statusUpdate struct {
		Status OrderStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		RespondError(w, http.StatusBadRequest, "Corps de requête invalide pour le statut") // ✅ Correction
		return
	}

	for i := range Orders {
		if Orders[i].ID == orderID {
			Orders[i].Status = statusUpdate.Status
			RespondJSON(w, http.StatusOK, Orders[i]) // ✅ Correction
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée") // ✅ Correction
}

// DeleteOrder : Handler exporté
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for i, order := range Orders {
		if order.ID == orderID {
			if order.Status == StatusPickedUp {
				RespondError(w, http.StatusBadRequest, "Impossible d'annuler une commande déjà récupérée") // ✅ Correction
				return
			}
			Orders = append(Orders[:i], Orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Commande non trouvée") // ✅ Correction
}
