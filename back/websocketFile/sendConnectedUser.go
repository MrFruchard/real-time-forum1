package websocketFile

import (
	"github.com/gorilla/websocket"
	"log"
	"real-time-forum/models"
)

func (h *Hub) sendConnectedUsers(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Vérifier si l'utilisateur est encore présent
	requesterSessionID, exists := h.clients[conn]
	if !exists {
		log.Println("L'utilisateur demandant la liste des connectés n'est plus dans le Hub.")
		err := conn.WriteJSON(models.UserStatus{
			Type:    "connected_users",
			Content: []string{"No users connected"},
		})
		if err != nil {
			log.Println("Erreur d'envoi du message 'No users connected':", err)
		}
		return
	}

	var users []string
	for _, sessionID := range h.clients {
		if sessionID != requesterSessionID {
			users = append(users, sessionID)
		}
	}

	// Si aucun autre utilisateur n'est connecté, renvoyer "No users connected"
	if len(users) == 0 {
		users = append(users, "No users connected")
	}

	msg := models.UserStatus{
		Type:    "connected_users",
		Content: users,
	}

	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("Erreur d'envoi de la liste des utilisateurs:", err)
	}
}
