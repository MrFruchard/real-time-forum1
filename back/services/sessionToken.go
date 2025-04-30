package services

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

// CreateSessionToken génère un token de session, supprime l'ancienne session si elle existe, puis enregistre la nouvelle.
func CreateSessionToken(db *sql.DB, userID string) (string, error) {
	// Début de la transaction
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'initialisation de la transaction : %w", err)
	}

	// Suppression de l'ancienne session
	deleteQuery := "DELETE FROM SESSION WHERE USERID = ?"
	_, err = tx.Exec(deleteQuery, userID)
	if err != nil {
		tx.Rollback() // Annule la transaction en cas d'erreur
		return "", fmt.Errorf("erreur lors de la suppression de l'ancienne session : %w", err)
	}

	// Génération d'un UUID pour la nouvelle session
	sessionToken, err := uuid.NewV4()
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("erreur lors de la génération du token de session : %w", err)
	}
	idSession := sessionToken.String()

	// Définition des timestamps
	now := time.Now()
	expiration := now.Add(24 * time.Hour) // Expiration dans 24 heures

	// Insertion de la nouvelle session
	insertQuery := "INSERT INTO SESSION (ID, USERID, CREATED_AT, LAST_ACTIVE_AT, EXPIRES_AT) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.Exec(insertQuery, idSession, userID, now, now, expiration)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("erreur lors de l'insertion en base de données : %w", err)
	}

	// Commit de la transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("erreur lors de la validation de la transaction : %w", err)
	}

	// Retourne le token de session
	return idSession, nil
}
