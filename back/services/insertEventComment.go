package services

import (
	"database/sql"
	"fmt"
	"log"
)

func CommentExist(db *sql.DB, postID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM COMMENT WHERE ID = ?)", postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la vérification de l'existence du post : %v", err)
	}
	return exists, nil
}

func CommentEventLike(db *sql.DB, userID string, commentID int) error {
	exists, err := CommentExist(db, commentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("le post avec ID %d n'existe pas", commentID)
	}
	// Vérifier si l'utilisateur a déjà liké
	var existingLike int
	err = db.QueryRow("SELECT COUNT(*) FROM LIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID).Scan(&existingLike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du like : %v", err)
	}

	if existingLike > 0 {
		// Supprimer le like s'il existe déjà
		_, err := db.Exec("DELETE FROM LIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du like : %v", err)
		}
		log.Println("Like supprimé")
		return nil
	}

	// Vérifier si un dislike existe
	var existingDislike int
	err = db.QueryRow("SELECT COUNT(*) FROM DISLIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID).Scan(&existingDislike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du dislike : %v", err)
	}

	if existingDislike > 0 {
		// Supprimer le dislike avant d'ajouter un like
		_, err := db.Exec("DELETE FROM DISLIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du dislike : %v", err)
		}
	}

	// Ajouter le like
	_, err = db.Exec("INSERT INTO LIKE_COMMENT (COMMENT_ID, USERID) VALUES (?, ?)", commentID, userID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du like : %v", err)
	}

	return nil
}

func CommentEventDislike(db *sql.DB, userID string, commentID int) error {
	exists, err := CommentExist(db, commentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("le post avec ID %d n'existe pas", commentID)
	}
	// Vérifier si l'utilisateur a déjà disliké
	var existingDislike int
	err = db.QueryRow("SELECT COUNT(*) FROM DISLIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID).Scan(&existingDislike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du dislike : %v", err)
	}

	if existingDislike > 0 {
		// Supprimer le dislike s'il existe déjà
		_, err := db.Exec("DELETE FROM DISLIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du dislike : %v", err)
		}
		log.Println("Dislike supprimé")
		return nil
	}

	// Vérifier si un like existe
	var existingLike int
	err = db.QueryRow("SELECT COUNT(*) FROM LIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID).Scan(&existingLike)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du like : %v", err)
	}

	if existingLike > 0 {
		// Supprimer le like avant d'ajouter un dislike
		_, err := db.Exec("DELETE FROM LIKE_COMMENT WHERE USERID = ? AND COMMENT_ID = ?", userID, commentID)
		if err != nil {
			return fmt.Errorf("erreur lors de la suppression du like : %v", err)
		}
	}

	// Ajouter le dislike
	_, err = db.Exec("INSERT INTO DISLIKE_COMMENT (COMMENT_ID, USERID) VALUES (?, ?)", commentID, userID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du dislike : %v", err)
	}

	return nil
}
