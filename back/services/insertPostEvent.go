package services

import (
	"database/sql"
	"fmt"
	"log"
)

func PostExists(db *sql.DB, postID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM POST WHERE ID = ?)", postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la vérification de l'existence du post : %v", err)
	}
	return exists, nil
}

// Gère l'ajout ou la suppression d'un like sur un post
func InsertPostLike(db *sql.DB, userID string, postID int) error {
	exists, err := PostExists(db, postID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("le post avec ID %d n'existe pas", postID)
	}
	// Vérifier si l'utilisateur a déjà liké
	var existingLike int
	err = db.QueryRow("SELECT COUNT(*) FROM LIKES WHERE USER_ID = ? AND POST_ID = ?", userID, postID).Scan(&existingLike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du like : %v", err)
	}

	if existingLike > 0 {
		// Supprimer le like s'il existe déjà
		_, err := db.Exec("DELETE FROM LIKES WHERE USER_ID = ? AND POST_ID = ?", userID, postID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du like : %v", err)
		}
		log.Println("Like supprimé")
		return nil
	}

	// Vérifier si un dislike existe
	var existingDislike int
	err = db.QueryRow("SELECT COUNT(*) FROM DISLIKE WHERE USER_ID = ? AND POST_ID = ?", userID, postID).Scan(&existingDislike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du dislike : %v", err)
	}

	if existingDislike > 0 {
		// Supprimer le dislike avant d'ajouter un like
		_, err := db.Exec("DELETE FROM DISLIKE WHERE USER_ID = ? AND POST_ID = ?", userID, postID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du dislike : %v", err)
		}
	}

	// Ajouter le like
	_, err = db.Exec("INSERT INTO LIKES (POST_ID, USER_ID) VALUES (?, ?)", postID, userID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du like : %v", err)
	}

	return nil
}

// Gère l'ajout ou la suppression d'un dislike sur un post
func InsertPostDislike(db *sql.DB, userID string, postID int) error {
	exists, err := PostExists(db, postID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("le post avec ID %d n'existe pas", postID)
	}
	// Vérifier si l'utilisateur a déjà disliké
	var existingDislike int
	err = db.QueryRow("SELECT COUNT(*) FROM DISLIKE WHERE USER_ID = ? AND POST_ID = ?", userID, postID).Scan(&existingDislike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du dislike : %v", err)
	}

	if existingDislike > 0 {
		// Supprimer le dislike s'il existe déjà
		_, err := db.Exec("DELETE FROM DISLIKE WHERE USER_ID = ? AND POST_ID = ?", userID, postID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du dislike : %v", err)
		}
		log.Println("Dislike supprimé")
		return nil
	}

	// Vérifier si un like existe
	var existingLike int
	err = db.QueryRow("SELECT COUNT(*) FROM LIKES WHERE USER_ID = ? AND POST_ID = ?", userID, postID).Scan(&existingLike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du like : %v", err)
	}

	if existingLike > 0 {
		// Supprimer le like avant d'ajouter un dislike
		_, err := db.Exec("DELETE FROM LIKES WHERE USER_ID = ? AND POST_ID = ?", userID, postID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du like : %v", err)
		}
	}

	// Ajouter le dislike
	_, err = db.Exec("INSERT INTO DISLIKE (POST_ID, USER_ID) VALUES (?, ?)", postID, userID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du dislike : %v", err)
	}

	return nil
}
