package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/utils"
)

type ReadNotification struct {
	Id int `json:"id"`
}

func Notification(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing cookie")
		return
	}
	checkSession, userID := utils.CheckSession(db, cookie.Value)
	if !checkSession {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Session invalid")
		return
	}

	if userID == "" {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Utilisateur invalide")
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleNotification(db, w, r, userID)
	case http.MethodPost:
		handleReadedNotification(db, w, r, userID)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleNotification(db *sql.DB, w http.ResponseWriter, r *http.Request, userId string) {
	notifications, err := utils.GetNotification(db, userId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Error lors de la recupération des données")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(notifications)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de l'encodage des données")
		return
	}
}

func handleReadedNotification(db *sql.DB, w http.ResponseWriter, r *http.Request, userId string) {
	var read ReadNotification
	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Error lors du décodage des données ")
	}

	err = utils.ReadNotification(db, userId, read.Id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Error lors de la récupération des données")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
