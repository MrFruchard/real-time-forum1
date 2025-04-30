package websocketFile

/*import (
	"database/sql"
	"encoding/json"
	"log"
	"real-time-forum/models"
	"real-time-forum/utils"
)

// SendConversation représente une conversation envoyée via WebSocket.
type SendConversation struct {
	Username2         string `json:"username2"`
	LastSenderMessage string `json:"lastSenderMessage"`
	LastMessage       string `json:"lastMessage"`
}

// BroadCastAllConversation envoie toutes les conversations de l'utilisateur via WebSocket.
func (h *Hub) BroadCastAllConversation(db *sql.DB, username string) {
	userID := utils.GetUserIDbyUsername(db, username)
	if userID == "" {
		log.Println("User ID not found for username:", username)
		return
	}

	query := `
		SELECT
			CASE
				WHEN USER1ID = ? THEN USER2ID
				ELSE USER1ID
			END AS other_user_id,
			(SELECT SENDER_ID FROM MESSAGE WHERE CONVERSATION_ID = c.ID ORDER BY SENT_AT DESC LIMIT 1) AS last_sender,
			(SELECT CONTENT FROM MESSAGE WHERE CONVERSATION_ID = c.ID ORDER BY SENT_AT DESC LIMIT 1) AS last_message
		FROM CONVERSATION c
		WHERE USER1ID = ? OR USER2ID = ?;
	`

	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		log.Println("❌ Error fetching conversations:", err)
		return
	}
	defer rows.Close()

	var conversations []string // Tableau de JSON en string

	for rows.Next() {
		var otherUserID, lastSenderID, lastMessage string

		err := rows.Scan(&otherUserID, &lastSenderID, &lastMessage)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		otherUsername := getusername(db, otherUserID)
		if otherUsername == "" {
			log.Println("Could not retrieve username for user ID:", otherUserID)
			continue
		}

		lastSenderUsername := getusername(db, lastSenderID)
		if lastSenderUsername == "" {
			log.Println("Could not retrieve username for sender ID:", lastSenderID)
			continue
		}

		// Créer une conversation sous forme JSON string
		conversationObj := SendConversation{
			Username2:         otherUsername,
			LastSenderMessage: lastSenderUsername,
			LastMessage:       lastMessage,
		}
		conversationJSON, err := json.Marshal(conversationObj)
		if err != nil {
			log.Println("Error marshalling conversation:", err)
			continue
		}

		conversations = append(conversations, string(conversationJSON)) // Ajouter la conversation en tant que JSON string
	}

	//Si l'utilisateur n'a pas de conversations, on envoie un tableau vide
	if len(conversations) == 0 {
		log.Println("No conversations found for user:", username)
		conversations = []string{}
	}

	// Création du message WebSocket avec `MessageConv`
	message := models.MessageConv{
		Type:    "conversation_list",
		Content: conversations, // Tableau de JSON en string
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for client, clientUsername := range h.clients {
		if clientUsername == username {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("Error sending conversation data:", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

// GetUsername récupère le nom d'utilisateur à partir de son ID
func getusername(db *sql.DB, userID string) string {
	var username string
	query := "SELECT USERNAME FROM USER WHERE ID = ? LIMIT 1"
	err := db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No username found for userID:", userID)
		} else {
			log.Println("Error retrieving username for userID:", userID, "-", err)
		}
		return ""
	}
	return username
}
*/
