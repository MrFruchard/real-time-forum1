package utils

import (
	"database/sql"
	"log"
)

// CheckSession vérifie si une session est valide et renvoie le userID associé.
func CheckSession(db *sql.DB, sessionId string) (bool, string) {
	if sessionId == "" {
		return false, ""
	}

	var userID string
	query := `SELECT USERID FROM SESSION WHERE ID = ? AND EXPIRES_AT > CURRENT_TIMESTAMP`
	err := db.QueryRow(query, sessionId).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Aucune session trouvée (expirée ou inexistante)
			log.Println("❌ Session invalide ou expirée:", sessionId)
			return false, ""
		}
		// Erreur inattendue
		log.Println("❌ Erreur de vérification de session:", err)
		return false, ""
	}

	// Session valide, renvoyer le userID
	return true, userID
}
