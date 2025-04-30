package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/services"
	"real-time-forum/utils"
	"strconv"
	"strings"
)

func Comment(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodPost:
		if userID == "" {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Utilisateur invalide")
			return
		}
		handleCreateComment(db, w, r, userID)
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/comment/")
		postID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
			return
		}
		handleComment(db, w, postID, userID)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleCreateComment(db *sql.DB, w http.ResponseWriter, r *http.Request, userId string) {
	var comment models.ReceiveComment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if comment.Content == "" || comment.IdPost == 0 {
		utils.SendErrorResponse(w, 400, "Missing required fields")
		return
	}

	err = services.CreateNotification(db, userId, "comment", comment.IdPost)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid Insertion DB")
		return
	}
	
	err = services.InsertComment(db, userId, comment.Content, comment.IdPost)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid Insertion DB")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleComment(db *sql.DB, w http.ResponseWriter, postID int, userID string) {
	// Récupérer les commentaires depuis la base de données
	comments, err := utils.GetCommentsByPostID(db, postID, userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error retrieving comments")
		return
	}

	// Retourner les commentaires sous forme de JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
