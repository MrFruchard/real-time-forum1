package utils

import (
	"database/sql"
	"errors"
	"log"
)

func ReadNotification(db *sql.DB, userId string, id int) error {
	// Vérifier si la notification appartient bien à l'utilisateur
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM NOTIFICATION 
			WHERE ID = ? AND RECEIVER_ID = ? AND STATUS = 'unread'
		)
	`, id, userId).Scan(&exists)
	if err != nil {
		log.Println("Erreur lors de la vérification de la notification :", err)
		return err
	}

	if !exists {
		return errors.New("notification non trouvée ou déjà lue")
	}

	// Mettre à jour le statut de la notification
	_, err = db.Exec(`
		UPDATE NOTIFICATION 
		SET STATUS = 'read' 
		WHERE ID = ? AND RECEIVER_ID = ?
	`, id, userId)
	if err != nil {
		log.Println("Erreur lors de la mise à jour de la notification :", err)
		return err
	}

	return nil
}
