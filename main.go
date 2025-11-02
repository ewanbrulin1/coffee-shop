package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// --- Helpers Génériques (Utilisés par tous les Handlers) ----------------------

// Fonction helper pour envoyer une réponse JSON
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erreur interne du serveur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Fonction helper pour envoyer une erreur HTTP
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

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
	// Créer le routeur Mux
	r := mux.NewRouter()

	// Route d'accueil (GET /)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Bienvenue à l'API du Coffee Shop !")
	}).Methods("GET")

	// Définition des Routes Menu (handlers dans drink.go)
	r.HandleFunc("/menu", getMenu).Methods("GET")
	r.HandleFunc("/menu/{id}", getDrink).Methods("GET")

	// Définition des Routes Commandes (handlers dans order.go)
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", updateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	fmt.Println("☕ Serveur de l'API Coffee Shop démarré sur http://localhost:8080")

	// Démarrer le serveur sur le port 8080 avec le middleware CORS
	err := http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %s\n", err)
	}
}
