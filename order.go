package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Order représente une commande
type Order struct {
	ID        string   `json:"id"`
	DrinkID   string   `json:"drink_id"`
	Extras    []string `json:"extras"`
	Status    string   `json:"status"` // "pending", "completed", "cancelled"
	Price     float64  `json:"price"`
	CreatedAt string   `json:"created_at"`
}

// Stockage en mémoire des commandes
var orders = []Order{}

// calculatePrice : Calcule le prix avec les extras
func calculatePrice(basePrice float64, extras []string) float64 {
	price := basePrice + (float64(len(extras)) * 0.5)
	return price
}

// CreateOrder : Crée une nouvelle commande avec extras
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DrinkID string   `json:"drink_id"`
		Extras  []string `json:"extras"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Trouve la boisson pour obtenir le prix de base
	var basePrice float64
	// Accède à 'drinks' depuis le package 'main' (défini dans drink.go)
	for _, drink := range drinks {
		if drink.ID == request.DrinkID {
			basePrice = drink.Price
			break
		}
	}

	// Calcule le prix final
	price := calculatePrice(basePrice, request.Extras)

	// Crée la commande
	newOrder := Order{
		ID:        fmt.Sprintf("%d", len(orders)+1),
		DrinkID:   request.DrinkID,
		Extras:    request.Extras,
		Status:    "pending",
		Price:     price,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	orders = append(orders, newOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

// GetOrders : Retourne la liste des commandes
func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus : Met à jour le statut d'une commande
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var update struct {
		Status string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = update.Status
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}

	http.NotFound(w, r)
}

// CancelOrder : Annule une commande
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = "cancelled"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(orders[i])
			return
		}
	}

	http.NotFound(w, r)
}
