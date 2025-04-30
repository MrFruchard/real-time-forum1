package utils

import (
	"database/sql"
	"fmt"
)

type UserConversation struct {
	ConversationID  string `json:"conversation_id"`
	Username        string `json:"username"`          // Nom de l'autre participant
	StartedAt       string `json:"started_at"`        // Date de début de la conversation
	LastMessage     string `json:"last_message"`      // Dernier message envoyé
	LastSender      string `json:"last_sender"`       // Nom de l'expéditeur du dernier message
	LastMessageDate string `json:"last_message_date"` // Date du dernier message
}

func GetUserConversationsWithLastMessage(db *sql.DB, userID string) ([]UserConversation, error) {
	query := `
	SELECT c.ID, u_other.USERNAME, c.STARTED_AT, 
	       COALESCE(m.CONTENT, '') AS last_message,
	       COALESCE(u_sender.USERNAME, '') AS last_sender,
	       COALESCE(m.SENT_AT, '') AS last_message_date
	FROM CONVERSATION c
	JOIN USER u_other ON (u_other.ID = c.USER1ID OR u_other.ID = c.USER2ID) 
	LEFT JOIN MESSAGE m ON m.ID = (
	    SELECT ID FROM MESSAGE 
	    WHERE MESSAGE.CONVERSATION_ID = c.ID 
	    ORDER BY MESSAGE.SENT_AT DESC 
	    LIMIT 1
	)
	LEFT JOIN USER u_sender ON u_sender.ID = m.SENDER_ID
	WHERE (c.USER1ID = ? OR c.USER2ID = ?) 
	AND u_other.ID != ? 
	ORDER BY c.STARTED_AT DESC;
	`

	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des conversations : %v", err)
	}
	defer rows.Close()

	var conversations []UserConversation

	for rows.Next() {
		var conv UserConversation
		if err := rows.Scan(&conv.ConversationID, &conv.Username, &conv.StartedAt, &conv.LastMessage, &conv.LastSender, &conv.LastMessageDate); err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture des conversations : %v", err)
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur après la récupération des conversations : %v", err)
	}

	return conversations, nil
}
