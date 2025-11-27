package model

// Drink : structure exportée
type Drink struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	BasePrice float64 `json:"base_price"`
}

// Drinks : données exportées
var Drinks = []Drink{
	{ID: "1", Name: "Espresso", Category: "coffee", BasePrice: 2.50},
	{ID: "2", Name: "Cappuccino", Category: "coffee", BasePrice: 3.50},
	{ID: "3", Name: "Latte", Category: "coffee", BasePrice: 4.00},
}
