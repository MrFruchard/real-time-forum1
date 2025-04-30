package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/models"
	"real-time-forum/utils"
	"strconv"
	"strings"
)

type UpdatePost struct {
	Content  string `json:"content"`
	Category string `json:"category"`
}

func UniquePost(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "ID du post manquant")
		return
	}

	// Convertir l'ID en entier
	postID, err := strconv.Atoi(pathParts[2])
	if err != nil || postID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "ID du post invalide")
		return
	}

	switch r.Method {
	case http.MethodPut:
		handleUpdatePost(db, w, r, postID, userID)
	case http.MethodDelete:
		handleDeletePost(db, w, postID, userID)
	case http.MethodGet:
		handleGetPost(db, w, postID, userID)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func handleDeletePost(db *sql.DB, w http.ResponseWriter, postID int, userID string) {
	// Vérifier si l'utilisateur est bien l'auteur du post
	var authorID string
	err := db.QueryRow("SELECT user_id FROM post WHERE id = ?", postID).Scan(&authorID)
	if err == sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusNotFound, "Post introuvable")
		return
	} else if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la vérification du post")
		return
	}

	// Vérifier que l'utilisateur a le droit de supprimer
	if authorID != userID {
		utils.SendErrorResponse(w, http.StatusForbidden, "Vous n'avez pas la permission de supprimer ce post")
		return
	}

	// Supprimer le post
	_, err = db.Exec("DELETE FROM post WHERE id = ?", postID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la suppression du post")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Post supprimé avec succès")
}

func handleUpdatePost(db *sql.DB, w http.ResponseWriter, r *http.Request, postID int, userID string) {
	// Vérifier si l'utilisateur est bien l'auteur du post
	var authorID string
	err := db.QueryRow("SELECT user_id FROM post WHERE id = ?", postID).Scan(&authorID)
	if err == sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusNotFound, "Post introuvable")
		return
	} else if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la vérification du post")
		return
	}

	// Vérifier que l'utilisateur a le droit de modifier
	if authorID != userID {
		utils.SendErrorResponse(w, http.StatusForbidden, "Vous n'avez pas la permission de modifier ce post")
		return
	}

	// Décoder le JSON de la requête
	var updateData UpdatePost
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Vérifier que le contenu et la catégorie ne sont pas vides
	if updateData.Content == "" || updateData.Category == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Le contenu et la catégorie sont obligatoires")
		return
	}

	// Mettre à jour le post dans la base de données
	_, err = db.Exec("UPDATE post SET content = ?, category = ? WHERE id = ?", updateData.Content, updateData.Category, postID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la mise à jour du post")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Post mis à jour avec succès")
}

func handleGetPost(db *sql.DB, w http.ResponseWriter, postID int, userID string) {
	// Définition de la structure de réponse
	var post models.ResponsePost

	err := db.QueryRow(`
		SELECT 
		    p.ID, 
		    u.USERNAME, 
		    p.CONTENT, 
		    p.IMAGE, 
		    p.CATEGORY, 
		    p.CREATED_AT,
		    COALESCE((SELECT COUNT(*) FROM LIKES WHERE POST_ID = p.ID), 0) AS Likes,
		    COALESCE((SELECT COUNT(*) FROM DISLIKE WHERE POST_ID = p.ID), 0) AS Dislikes,
		    COALESCE((SELECT COUNT(*) FROM COMMENT WHERE POST_ID = p.ID), 0) AS Comments,
		    EXISTS (SELECT 1 FROM LIKES WHERE POST_ID = p.ID AND USER_ID = ?) AS Liked,
		    EXISTS (SELECT 1 FROM DISLIKE WHERE POST_ID = p.ID AND USER_ID = ?) AS Disliked
		FROM POST p
		JOIN USER u ON p.USER_ID = u.ID
		WHERE p.ID = ?
	`, userID, userID, postID).
		Scan(&post.Id, &post.User, &post.Content, &post.Image, &post.Category, &post.CreatedAt,
			&post.Likes, &post.Dislikes, &post.Comments, &post.Liked, &post.Disliked)

	// Gestion des erreurs
	if err == sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusNotFound, "Post introuvable")
		return
	} else if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération du post")
		return
	}

	// Envoyer la réponse JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
