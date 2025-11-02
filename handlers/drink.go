package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ⚠️ Plus d'importation problématique du module racine.

// Drink : structure exportée
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	BasePrice float64 `json:"base_price"`
}

// Drinks : données exportées
var Drinks []Drink = []Drink{
	{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
	{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
	{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 4.00},
	{ID: "4", Name: "Americano", Category: "coffee", BasePrice: 3.00},
	{ID: "5", Name: "Green Tea", Category: "tea", BasePrice: 2.50},
	{ID: "6", Name: "Iced Coffee", Category: "cold", BasePrice: 4.50},
}

// GetMenu : Handler exporté (GET /menu)
func GetMenu(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, Drinks) // ✅ Correction
}

// GetDrink : Handler exporté (GET /menu/{id})
func GetDrink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	drinkID := vars["id"]

	for _, drink := range Drinks {
		if drink.ID == drinkID {
			RespondJSON(w, http.StatusOK, drink) // ✅ Correction
			return
		}
	}

	RespondError(w, http.StatusNotFound, "Boisson non trouvée") // ✅ Correction
}

// FindDrinkByID : Helper exporté (utilisé par order.go)
func FindDrinkByID(id string) (Drink, bool) {
	for _, drink := range Drinks {
		if drink.ID == id {
			return drink, true
		}
	}
	return Drink{}, false
}
