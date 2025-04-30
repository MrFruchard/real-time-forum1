package websocketFile

import (
	"github.com/gorilla/websocket"
	"log"
	"real-time-forum/models"
)

func (h *Hub) BroadcastDisconnectedUser(username string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	msg := models.UserStatus{
		Type:    "user_disconnected",
		Content: []string{username},
	}

	// Liste des clients à supprimer
	toRemove := make([]*websocket.Conn, 0)

	for client, clientSessionID := range h.clients {
		// ✅ Ne pas envoyer à l'utilisateur qui se déconnecte
		if clientSessionID == username {
			toRemove = append(toRemove, client) // Ajouter pour suppression
			continue
		}

		// Envoyer aux autres utilisateurs
		err := client.WriteJSON(msg)
		if err != nil {
			// Vérifier si c'est une fermeture propre
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("ℹ️ Client fermé proprement :", err)
			} else {
				log.Println("❌ Erreur d'envoi du message de déconnexion :", err)
			}
			toRemove = append(toRemove, client) // Ajouter pour suppression
		}
	}

	// Supprimer directement les connexions fermées sans relâcher le verrou
	for _, client := range toRemove {
		delete(h.clients, client)
	}
}
