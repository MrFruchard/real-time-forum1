package utils

import (
	"database/sql"
	"real-time-forum/models"
)

func GetAllMessage(db *sql.DB, convId string) ([]models.SendAllPrivateMessage, error) {
	// Définition de la requête SQL
	query := `
		SELECT SENDER_ID, CONTENT, SENT_AT 
		FROM MESSAGE 
		WHERE CONVERSATION_ID = ? 
		ORDER BY SENT_AT ASC
	`

	// Exécuter la requête
	rows, err := db.Query(query, convId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Stocker les résultats
	var messages []models.SendAllPrivateMessage

	for rows.Next() {
		var senderID, content, sentAt string

		if err := rows.Scan(&senderID, &content, &sentAt); err != nil {
			return nil, err
		}

		// Récupérer le username depuis la table USER
		username, err := GetUsernamebyID(db, senderID)
		if err != nil {
			return nil, err
		}

		// Ajouter le message avec le nom d'utilisateur
		messages = append(messages, models.SendAllPrivateMessage{
			Sender:  username,
			Message: content,
			Date:    sentAt,
		})
	}

	// Vérifier les erreurs après l'itération
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// Fonction pour récupérer le username via l'ID utilisateur
func GetUsernamebyID(db *sql.DB, userID string) (string, error) {
	var username string
	err := db.QueryRow("SELECT USERNAME FROM USER WHERE ID = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
