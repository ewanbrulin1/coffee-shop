package handlers

import (
	"encoding/json"
	"net/http"
)

// RespondJSON : Fonction helper pour envoyer une réponse JSON
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erreur interne du serveur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondError : Fonction helper pour envoyer une erreur HTTP
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}
