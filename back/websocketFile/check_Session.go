package websocketFile

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func isValidSession(db *sql.DB, sessionID string) bool {
	var userID string
	var expiresAtStr string

	err := db.QueryRow("SELECT USERID, EXPIRES_AT FROM SESSION WHERE ID = ?", sessionID).Scan(&userID, &expiresAtStr)
	if err != nil {
		log.Println("Erreur SQL:", err)
		return false
	}

	// Parser la date selon le format SQLite
	expiresAt, err := time.Parse("2006-01-02 15:04:05.999999-07:00", expiresAtStr)
	if err != nil {
		// Essayer un autre format si l'erreur persiste
		expiresAt, err = time.Parse("2006-01-02 15:04:05", expiresAtStr)
		if err != nil {
			log.Println("Erreur de parsing de la date:", err)
			return false
		}
	}

	// Vérifier si la session est expirée
	if time.Now().UTC().After(expiresAt.UTC()) {
		fmt.Println("Session expirée")
		return false
	}

	fmt.Println("Session valide")
	return true // Session valide et non expirée
}
