package handlers

import (
	"database/sql"
	"net/http"
	"real-time-forum/utils"
	"time"
)

func Logout(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing cookie")
		return
	}
	sessionId := cookie.Value

	checkSession, _ := utils.CheckSession(db, cookie.Value)
	if !checkSession {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Session invalid")
		return
	}
	switch r.Method {
	case http.MethodPost:
		handleLogout(db, w, r, sessionId)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Methode non autorisé")
	}
}

func handleLogout(db *sql.DB, w http.ResponseWriter, r *http.Request, sessionId string) {
	_, err := db.Exec("DELETE FROM SESSION WHERE ID = ?", sessionId)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression de la session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Date passée pour forcer l'expiration
		MaxAge:   -1,              // Supprime immédiatement le cookie
		HttpOnly: true,
		Secure:   true, // À activer en HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}
