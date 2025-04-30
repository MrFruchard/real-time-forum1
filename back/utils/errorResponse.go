package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json") // Définit les headers avant d'écrire
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(ErrorMessage{
		Code:    code,
		Message: message,
	})

	if err != nil {
		log.Printf("Erreur lors de l'encodage JSON : %v", err)
	}
}

type SuccessMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendSuccessResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json") // Définit les headers avant d'écrire
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(SuccessMessage{
		Code:    code,
		Message: message,
	})

	if err != nil {
		log.Printf("Erreur lors de l'encodage JSON : %v", err)
	}
}
