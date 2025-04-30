package services

import (
	"database/sql"
	"errors"
	"log"
)

func CreateNotification(db *sql.DB, userID, typeNotif string, id int) error {
	var receiverID string

	// Déterminer le receiverID en fonction du type de notification
	switch typeNotif {
	case "comment", "like", "dislike":
		if err := db.QueryRow(`SELECT USER_ID FROM POST WHERE ID = ?`, id).Scan(&receiverID); err != nil {
			if err == sql.ErrNoRows {
				return errors.New("post non trouvé")
			}
			log.Println("Erreur lors de la récupération du destinataire pour un post:", err)
			return err
		}
	case "like_comment", "dislike_comment":
		if err := db.QueryRow(`SELECT USERID FROM COMMENT WHERE ID = ?`, id).Scan(&receiverID); err != nil {
			if err == sql.ErrNoRows {
				return errors.New("commentaire non trouvé")
			}
			log.Println("Erreur lors de la récupération du destinataire pour un commentaire:", err)
			return err
		}
	default:
		log.Println("Type de notification inconnu:", typeNotif)
		return errors.New("type de notification invalide")
	}

	// Ne pas créer de notification pour soi-même
	if receiverID == userID {
		return nil
	}

	// Gestion spécifique pour les notifications de type "comment" : on les empile
	if typeNotif == "comment" {
		_, err := db.Exec(`
			INSERT INTO NOTIFICATION (RECEIVER_ID, SENDER_ID, TYPE, RELATED_ID, STATUS, CREATED_AT)
			VALUES (?, ?, ?, ?, 'unread', datetime('now', 'localtime'))
		`, receiverID, userID, typeNotif, id)
		if err != nil {
			log.Println("Erreur lors de la création de la notification (comment) :", err)
		}
		return err
	}

	// Pour les notifications "toggleables" (like/dislike ou like_comment/dislike_comment),
	// on définit le groupe de types concernés.
	var toggleTypes [2]string
	switch typeNotif {
	case "like", "dislike":
		toggleTypes = [2]string{"like", "dislike"}
	case "like_comment", "dislike_comment":
		toggleTypes = [2]string{"like_comment", "dislike_comment"}
	default:
		// Pour d'autres types, on se contente d'insérer
		_, err := db.Exec(`
			INSERT INTO NOTIFICATION (RECEIVER_ID, SENDER_ID, TYPE, RELATED_ID, STATUS, CREATED_AT)
			VALUES (?, ?, ?, ?, 'unread', datetime('now', 'localtime'))
		`, receiverID, userID, typeNotif, id)
		return err
	}

	// Vérifier si une notification du même groupe (toggle) existe déjà pour ce (receiver, sender, id)
	query := `
		SELECT TYPE FROM NOTIFICATION 
		WHERE RECEIVER_ID = ? AND SENDER_ID = ? AND RELATED_ID = ? 
		AND TYPE IN (?, ?)
	`
	var existingType string
	err := db.QueryRow(query, receiverID, userID, id, toggleTypes[0], toggleTypes[1]).Scan(&existingType)
	if err == nil {
		// Une notification de like/dislike existe déjà dans ce groupe
		if existingType == typeNotif {
			// Même action effectuée une seconde fois : on supprime la notification (toggle off)
			delQuery := `
				DELETE FROM NOTIFICATION 
				WHERE RECEIVER_ID = ? AND SENDER_ID = ? AND RELATED_ID = ? AND TYPE = ?
			`
			if _, err := db.Exec(delQuery, receiverID, userID, id, typeNotif); err != nil {
				log.Println("Erreur lors de la suppression de la notification :", err)
				return err
			}
			return nil
		} else {
			// Action différente dans le même groupe (ex: like → dislike)
			// On supprime l'ancienne notification puis on continue pour insérer la nouvelle.
			delQuery := `
				DELETE FROM NOTIFICATION 
				WHERE RECEIVER_ID = ? AND SENDER_ID = ? AND RELATED_ID = ? 
				AND TYPE IN (?, ?)
			`
			if _, err := db.Exec(delQuery, receiverID, userID, id, toggleTypes[0], toggleTypes[1]); err != nil {
				log.Println("Erreur lors de la suppression de l'ancienne notification :", err)
				return err
			}
		}
	} else if err != sql.ErrNoRows {
		log.Println("Erreur lors de la vérification des notifications existantes :", err)
		return err
	}

	// Insérer la nouvelle notification toggleable (like/dislike ou like_comment/dislike_comment)
	insertQuery := `
		INSERT INTO NOTIFICATION (RECEIVER_ID, SENDER_ID, TYPE, RELATED_ID, STATUS, CREATED_AT)
		VALUES (?, ?, ?, ?, 'unread', datetime('now', 'localtime'))
	`
	if _, err := db.Exec(insertQuery, receiverID, userID, typeNotif, id); err != nil {
		log.Println("Erreur lors de la création de la notification :", err)
		return err
	}

	return nil
}
