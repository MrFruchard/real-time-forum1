package services

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"log"
)

func InsertNewConv(db *sql.DB, senderId, receiverId string) (string, error) {
	var convId string

	// Vérifier si une conversation existe déjà entre sender et receiver (dans les deux sens)
	query := `
		SELECT ID FROM CONVERSATION 
		WHERE (USER1ID = ? AND USER2ID = ?) OR (USER1ID = ? AND USER2ID = ?)
	`
	err := db.QueryRow(query, senderId, receiverId, receiverId, senderId).Scan(&convId)

	if err == nil {
		// Si une conversation existe déjà, on retourne son ID
		log.Println("Conversation existante trouvée avec ID :", convId)
		return convId, nil
	} else if err != sql.ErrNoRows {
		// S'il y a une autre erreur SQL, on la log
		log.Println("Erreur lors de la vérification de la conversation :", err)
		return "", err
	}

	log.Println("Aucune conversation existante trouvée entre", senderId, "et", receiverId, "→ création d'une nouvelle conversation...")

	// Insérer une nouvelle conversation avec senderId en USER1ID et receiverId en USER2ID
	insertQuery := `INSERT INTO CONVERSATION (ID, USER1ID, USER2ID) VALUES (?, ?, ?)`
	convUUID, _ := uuid.NewV4()
	convId = convUUID.String()

	_, err = db.Exec(insertQuery, convId, senderId, receiverId)
	if err != nil {
		log.Println("Erreur lors de l'insertion de la conversation :", err)
		return "", err
	}

	log.Println("Nouvelle conversation créée avec ID :", convId)
	return convId, nil
}
