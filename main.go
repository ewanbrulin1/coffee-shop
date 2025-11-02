package main

import (
	"fmt"
	"net/http"

	"github.com/ewanb/CoffeeGo/coffee-shop-api/handlers" // Importation du package handlers
	"github.com/gorilla/mux"
)

// NOTE: Les fonctions RespondJSON et RespondError sont maintenant dans handlers/utils.go

// --- Middleware CORS ---------------------------------------------------------

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --- Fonction Main -----------------------------------------------------------

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Bienvenue à l'API du Coffee Shop !")
	}).Methods("GET")

	// Routes Menu
	r.HandleFunc("/menu", handlers.GetMenu).Methods("GET")
	r.HandleFunc("/menu/{id}", handlers.GetDrink).Methods("GET")

	// Routes Commandes
	r.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/orders", handlers.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", handlers.UpdateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", handlers.DeleteOrder).Methods("DELETE")

	fmt.Println("☕ Serveur de l'API Coffee Shop démarré sur http://localhost:8080")

	err := http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %s\n", err)
	}
}
