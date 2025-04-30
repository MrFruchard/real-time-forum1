package utils

import (
	"database/sql"
	"log"
	"real-time-forum/models"
)

// GetCommentsByPostID récupère tous les commentaires associés à un post donné, avec les likes et dislikes de l'utilisateur.
func GetCommentsByPostID(db *sql.DB, postID int, userID string) ([]models.ResponseComment, error) {
	var comments []models.ResponseComment

	rows, err := db.Query(`
		SELECT c.ID, u.USERNAME, c.CONTENT, 
		       COALESCE(l.Likes, 0) AS Likes, 
		       COALESCE(d.Dislikes, 0) AS Dislikes, 
		       c.CREATED_AT,
		       CASE WHEN ul.USERID IS NOT NULL THEN 1 ELSE 0 END AS Liked,
		       CASE WHEN ud.USERID IS NOT NULL THEN 1 ELSE 0 END AS Disliked
		FROM COMMENT c
		JOIN USER u ON c.USERID = u.ID
		LEFT JOIN (SELECT COMMENT_ID, COUNT(*) AS Likes FROM LIKE_COMMENT GROUP BY COMMENT_ID) l ON c.ID = l.COMMENT_ID
		LEFT JOIN (SELECT COMMENT_ID, COUNT(*) AS Dislikes FROM DISLIKE_COMMENT GROUP BY COMMENT_ID) d ON c.ID = d.COMMENT_ID
		LEFT JOIN LIKE_COMMENT ul ON c.ID = ul.COMMENT_ID AND ul.USERID = ?
		LEFT JOIN DISLIKE_COMMENT ud ON c.ID = ud.COMMENT_ID AND ud.USERID = ?
		WHERE c.POST_ID = ?
		ORDER BY c.CREATED_AT ASC
	`, userID, userID, postID)
	if err != nil {
		log.Println("Erreur lors de la récupération des commentaires pour le post", postID, ":", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.ResponseComment
		if err := rows.Scan(&comment.Id, &comment.User, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.CreatedAt, &comment.Liked, &comment.Disliked); err != nil {
			log.Println("Erreur lors du scan des commentaires:", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		log.Println("Erreur lors de l'itération des commentaires:", err)
		return nil, err
	}

	return comments, nil
}
