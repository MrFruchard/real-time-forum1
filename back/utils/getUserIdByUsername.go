package utils

import (
	"database/sql"
	"log"
)

func GetUserIDbyUsername(db *sql.DB, username string) string {
	var userID string
	query := `SELECT ID FROM USER WHERE USERNAME = ?`
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Aucun utilisateur trouvé avec le nom d'utilisateur : %s", username)
		} else {
			log.Printf("Erreur lors de la récupération de l'ID utilisateur : %v", err)
		}
		return ""
	}
	return userID
}
