package services

import (
	"database/sql"
	"fmt"
	"time"
)

// InsertPost ajoute un post dans la base de données avec la date de création.
func InsertPost(db *sql.DB, userID, content, category, image string) error {
	query := `INSERT INTO POST (USER_ID, CONTENT, CATEGORY, IMAGE, CREATED_AT) 
	          VALUES (?, ?, ?, ?, ?)` // ✅ Correction de la requête SQL

	createdAt := time.Now() // ✅ Ajout de la date actuelle
	image = "hello"
	_, err := db.Exec(query, userID, content, category, image, createdAt)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion en base de données : %w", err)
	}

	return nil
}
