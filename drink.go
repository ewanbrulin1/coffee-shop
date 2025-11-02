package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Drink représente une boisson du menu
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"` // coffee, tea, cold
	BasePrice float64 `json:"base_price"`
}

// Base de données en mémoire pour le menu
var drinks []Drink = []Drink{
	{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
	{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
	{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 4.00},
	{ID: "4", Name: "Americano", Category: "coffee", BasePrice: 3.00},
	{ID: "5", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
	{ID: "6", Name: "Iced Coffee", Category: "cold", BasePrice: 4.50},
}

// Handler pour GET /menu - Récupérer le menu complet
func getMenu(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, drinks)
}

// Handler pour GET /menu/{id} - Récupérer une boisson spécifique
func getDrink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	drinkID := vars["id"]

	for _, drink := range drinks {
		if drink.ID == drinkID {
			respondJSON(w, http.StatusOK, drink)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Boisson non trouvée")
}

// Helper interne pour trouver une boisson (utilisé par order.go)
func findDrinkByID(id string) (Drink, bool) {
	for _, drink := range drinks {
		if drink.ID == id {
			return drink, true
		}
	}
	return Drink{}, false
}
