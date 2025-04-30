package utils

import (
	"database/sql"
	"log"
	"real-time-forum/models"
)

// GetPostsWithComments récupère tous les posts avec leur nombre de commentaires, likes, dislikes et si l'utilisateur a aimé ou non.
func GetPostsWithComments(db *sql.DB, userID string) ([]models.ResponsePost, error) {
	var posts []models.ResponsePost

	rows, err := db.Query(`
		SELECT p.ID, u.USERNAME, p.CONTENT, p.IMAGE, p.CATEGORY, p.CREATED_AT,
		       COALESCE(l.Likes, 0) AS Likes, 
		       COALESCE(d.Dislikes, 0) AS Dislikes,
		       COALESCE(c.CommentCount, 0) AS CommentCount,
		       CASE WHEN ul.USER_ID IS NOT NULL THEN 1 ELSE 0 END AS Liked,
		       CASE WHEN ud.USER_ID IS NOT NULL THEN 1 ELSE 0 END AS Disliked
		FROM POST p
		JOIN USER u ON p.USER_ID = u.ID
		LEFT JOIN (SELECT POST_ID, COUNT(*) AS Likes FROM LIKES GROUP BY POST_ID) l ON p.ID = l.POST_ID
		LEFT JOIN (SELECT POST_ID, COUNT(*) AS Dislikes FROM DISLIKE GROUP BY POST_ID) d ON p.ID = d.POST_ID
		LEFT JOIN (SELECT POST_ID, COUNT(*) AS CommentCount FROM COMMENT GROUP BY POST_ID) c ON p.ID = c.POST_ID
		LEFT JOIN LIKES ul ON p.ID = ul.POST_ID AND ul.USER_ID = ?
		LEFT JOIN DISLIKE ud ON p.ID = ud.POST_ID AND ud.USER_ID = ?
		ORDER BY p.CREATED_AT DESC
	`, userID, userID)

	if err != nil {
		log.Println("Erreur lors de la récupération des posts:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.ResponsePost
		if err := rows.Scan(&post.Id, &post.User, &post.Content, &post.Image, &post.Category, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments, &post.Liked, &post.Disliked); err != nil {
			log.Println("Erreur lors du scan des posts:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Println("Erreur lors de l'itération des posts:", err)
		return nil, err
	}

	return posts, nil
}
