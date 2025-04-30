package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/services"
	"real-time-forum/utils"
)

func Post(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	case http.MethodPost:
		handleCreatePost(w, r, db, userID)
	case http.MethodGet:
		handleGetPosts(w, r, db, userID)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

func handleCreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Vérification des champs obligatoires
	if post.Content == "" || post.Category == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Le contenu et la catégorie sont obligatoires")
		return
	}

	err = services.InsertPost(db, userID, post.Content, post.Category, post.Image)
	if err != nil {
		log.Printf("Erreur lors de l'insertion du post : %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de l'insertion du post")
		return
	}

	// Réponse JSON propre
	response := map[string]string{
		"status":  "Created",
		"message": "Post inséré avec succès",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleGetPosts(w http.ResponseWriter, r *http.Request, db *sql.DB, userID string) {
	posts, err := utils.GetPostsWithComments(db, userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des posts")
		return
	}

	// Vérification de l'encodage JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Erreur d'encodage JSON : %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de l'envoi des posts")
	}
}
