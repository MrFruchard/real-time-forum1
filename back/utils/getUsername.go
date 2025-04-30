package utils

import (
	"database/sql"
	"log"
)

func GetUsername(db *sql.DB, sessionId string) string {
	var username string
	query := `
		SELECT u.USERNAME
		FROM USER u
		JOIN SESSION s ON u.ID = s.USERID
		WHERE s.ID = ?`
	err := db.QueryRow(query, sessionId).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Session ID not found")
			return ""
		}
		log.Println("Error fetching username:", err)
		return ""
	}
	return username
}
