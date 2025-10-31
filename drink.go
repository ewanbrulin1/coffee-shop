package main

import (
	"encoding/json"
	"net/http"
)

// Drink représente une boisson
type Drink struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Stockage en mémoire des boissons
var drinks = []Drink{
	{ID: "1", Name: "Espresso", Price: 2.5},
	{ID: "2", Name: "Cappuccino", Price: 3.5},
	{ID: "3", Name: "Latte", Price: 3.0},
}

// getMenu : Retourne la liste des boissons
func getMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drinks)
}

// GetDrinkNameByID : Retourne le nom d'une boisson à partir de son ID
func GetDrinkNameByID(id string) string {
	for _, drink := range drinks {
		if drink.ID == id {
			return drink.Name
		}
	}
	return "Inconnu"
}
