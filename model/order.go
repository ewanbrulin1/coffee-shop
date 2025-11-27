package model

import "time"

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

// CalculatePrice : fonction exportée
func CalculatePrice(basePrice float64, size string, extras []string) float64 {
	price := basePrice
	switch size {
	case "small":
		price *= 0.8
	case "large":
		price *= 1.3
	}
	price += float64(len(extras)) * 0.50
	return price
}
