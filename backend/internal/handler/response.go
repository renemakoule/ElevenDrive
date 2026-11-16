package handler

import (
	"encoding/json"
	"net/http"
)

// StandardResponse structure standard de toutes les réponses JSON de l'API.
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON écrit une réponse JSON formatée avec le code HTTP spécifié.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success renvoie une réponse HTTP 200/201 contenant les données sérialisées.
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	JSON(w, statusCode, StandardResponse{
		Success: true,
		Data:    data,
	})
}

// Error renvoie une réponse d'erreur HTTP avec le message explicatif.
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, StandardResponse{
		Success: false,
		Error:   message,
	})
}
