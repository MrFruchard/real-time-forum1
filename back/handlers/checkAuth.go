// handlers/checkAuth.go
package handlers

import (
	"database/sql"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Vérifier la validité de la session dans la base de données
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM SESSION WHERE ID = ? AND EXPIRES_AT > CURRENT_TIMESTAMP)`
	err = db.QueryRow(query, cookie.Value).Scan(&exists)

	if err != nil || !exists {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	// Si tout est ok, renvoyer 200
	w.WriteHeader(http.StatusOK)
}
