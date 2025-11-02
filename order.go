package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// OrderStatus représente l'état d'une commande
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked-up"
	StatusCancelled OrderStatus = "cancelled"
)

// Order représente une commande
type Order struct {
	ID           string      `json:"id"`
	DrinkID      string      `json:"drink_id"`
	DrinkName    string      `json:"drink_name"`
	Size         string      `json:"size"` // small, medium, large
	Extras       []string    `json:"extras"`
	CustomerName string      `json:"customer_name"`
	Status       OrderStatus `json:"status"`
	TotalPrice   float64     `json:"total_price"`
	OrderedAt    time.Time   `json:"ordered_at"`
}

// Base de données en mémoire pour les commandes
var orders []Order
var orderCounter int = 1

// calculatePrice : Fonction helper pour calculer le prix total
func calculatePrice(basePrice float64, size string, extras []string) float64 {
	price := basePrice

	// Ajustement tailles
	switch size {
	case "small":
		price *= 0.8
	case "medium":
		// prix * 1.0
	case "large":
		price *= 1.3
	}

	// Ajout extras (0.50€ par extra)
	price += float64(len(extras)) * 0.50

	return price
}

// Handler pour POST /orders - Créer une nouvelle commande
func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		respondError(w, http.StatusBadRequest, "Corps de requête invalide")
		return
	}

	// Vérifier que la boisson (DrinkID) existe (utilise findDrinkByID de drink.go)
	foundDrink, drinkFound := findDrinkByID(newOrder.DrinkID)
	if !drinkFound {
		respondError(w, http.StatusBadRequest, "ID de boisson introuvable")
		return
	}

	// Générer un ID unique et remplir les champs
	newOrder.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	orderCounter++
	newOrder.DrinkName = foundDrink.Name
	newOrder.Status = StatusPending
	newOrder.OrderedAt = time.Now()
	// Gère les extras et la taille dans le calcul du prix
	newOrder.TotalPrice = calculatePrice(foundDrink.BasePrice, newOrder.Size, newOrder.Extras)

	// Ajouter la commande au slice orders
	orders = append(orders, newOrder)

	respondJSON(w, http.StatusCreated, newOrder)
}

// Handler pour GET /orders - Récupérer toutes les commandes
func getOrders(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, orders)
}

// Handler pour GET /orders/{id} - Récupérer une commande spécifique
func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for _, order := range orders {
		if order.ID == orderID {
			respondJSON(w, http.StatusOK, order)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Commande non trouvée")
}

// Handler pour PATCH /orders/{id}/status - Changer le statut d'une commande
func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	// Struct temporaire pour décoder uniquement le statut
	var statusUpdate struct {
		Status OrderStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		respondError(w, http.StatusBadRequest, "Corps de requête invalide pour le statut")
		return
	}

	for i := range orders {
		if orders[i].ID == orderID {
			orders[i].Status = statusUpdate.Status
			respondJSON(w, http.StatusOK, orders[i])
			return
		}
	}

	respondError(w, http.StatusNotFound, "Commande non trouvée")
}

// Handler pour DELETE /orders/{id} - Annuler une commande
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	for i, order := range orders {
		if order.ID == orderID {
			// Vérifier que le statut n'est pas "picked-up"
			if order.Status == StatusPickedUp {
				respondError(w, http.StatusBadRequest, "Impossible d'annuler une commande déjà récupérée")
				return
			}
			// Supprimer la commande du slice
			orders = append(orders[:i], orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Commande non trouvée")
}
