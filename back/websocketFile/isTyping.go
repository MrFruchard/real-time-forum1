package websocketFile

import (
	"log"
	"real-time-forum/models"
)

func (h *Hub) IsTyping(targetUsername string, senderUsername string, typing bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	msg := models.IsTyping{
		Type:    "is_typing",
		Content: senderUsername,
		Typing:  typing,
	}

	// Trouver la connexion du destinataire
	for client, username := range h.clients {
		if username == targetUsername {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Erreur lors de l'envoi du message:", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
