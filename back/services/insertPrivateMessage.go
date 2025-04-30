package services

import (
	"database/sql"
	"log"
)

func InsertNewMessage(db *sql.DB, conversationId string, senderId, content string) (int64, error) {
	// Requête d'insertion
	insertQuery := `INSERT INTO MESSAGE (CONVERSATION_ID, SENDER_ID, CONTENT) VALUES (?, ?, ?)`

	// Exécuter la requête
	result, err := db.Exec(insertQuery, conversationId, senderId, content)
	if err != nil {
		log.Println("Erreur lors de l'insertion du message :", err)
		return 0, err
	}

	// Récupérer l'ID du message inséré
	messageId, err := result.LastInsertId()
	if err != nil {
		log.Println("Erreur lors de la récupération de l'ID du message :", err)
		return 0, err
	}

	return messageId, nil
}
