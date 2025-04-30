package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/services"
	"real-time-forum/utils"
	"time"

	"github.com/gofrs/uuid"
)

type RegisterInsert struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       string `json:"age"`
	Genre     string `json:"genre"`
}

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if db == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur : base de données non disponible")
		return
	}

	var user RegisterInsert
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Erreur lors du décodage du JSON")
		return
	}
	defer r.Body.Close()

	// Vérification des champs obligatoires
	if user.Email == "" || user.Password == "" || user.Username == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Les champs email, mot de passe et nom d'utilisateur sont obligatoires")
		return
	}

	// Génération de l'UUID
	u, err := uuid.NewV4()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la génération de l'UUID")
		return
	}
	user.ID = u.String()

	// Canaux pour la goroutine
	hashedPassChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(hashedPassChan)
		defer close(errChan)

		hashedPass, err := utils.HashPassword(user.Password)
		if err != nil {
			errChan <- err
			return
		}
		hashedPassChan <- hashedPass
	}()

	if !checkDataRegister(db, user.Email, user.Username) {
		utils.SendErrorResponse(w, http.StatusConflict, "L'email ou le nom d'utilisateur existe déjà")
		return
	}

	var hashedPass string
	select {
	case hashedPass = <-hashedPassChan:
	case err = <-errChan:
		fmt.Printf("Erreur lors du hashage : %v\n", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Erreur lors du traitement interne")
		return
	}

	err = services.InsertUser(db, user.ID, user.Email, hashedPass, user.Username, user.FirstName, user.LastName, user.Age, user.Genre)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionToken, err := services.CreateSessionToken(db, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Définir un cookie HTTP-Only pour stocker la session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,  // Empêche l'accès via JavaScript (protection XSS)
		Secure:   false, // Mettre true en production (HTTPS obligatoire)
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Réponse JSON de succès
	response := map[string]string{
		"message": "Inscription réussie",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func checkDataRegister(db *sql.DB, email, username string) bool {
	if !utils.IsvalidEmail(email) {
		fmt.Println("Email invalide")
		return false
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM USER WHERE EMAIL = ? OR USERNAME = ?)`
	err := db.QueryRow(query, email, username).Scan(&exists)
	if err != nil {
		fmt.Printf("Erreur lors de la vérification des données : %v\n", err)
		return false
	}
	return !exists
}
