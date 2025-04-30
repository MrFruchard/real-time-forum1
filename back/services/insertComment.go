package services

import (
	"database/sql"
	"log"
)

// InsertComment insère un commentaire dans la base de données.
func InsertComment(db *sql.DB, userId, content string, postId int) error {
	query := `INSERT INTO Comment (USERID, CONTENT, POST_ID, CREATED_AT) VALUES (?, ?, ?, datetime('now'))`

	_, err := db.Exec(query, userId, content, postId)
	if err != nil {
		log.Printf("Error inserting comment: %v", err)
		return err
	}

	return nil
}
