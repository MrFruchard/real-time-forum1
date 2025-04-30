package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/utils"
)

func Conversation(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	sendConversation(db, w, r, userID)
}

func sendConversation(db *sql.DB, w http.ResponseWriter, r *http.Request, userId string) {
	// Récupération des conversations
	conversations, err := utils.GetUserConversationsWithLastMessage(db, userId)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des conversations", http.StatusInternalServerError)
		return
	}

	// Convertir en JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversations)
}
