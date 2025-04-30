package websocketFile

import (
	"github.com/gorilla/websocket"
	"log"
	"real-time-forum/models"
)

func (h *Hub) broadcastNewUser(username string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("✅ Nouvel utilisateur connecté : %s\n", username)

	msg := models.UserStatus{
		Type:    "new_user",
		Content: []string{username},
	}

	// Liste des clients à supprimer après la diffusion du message
	toRemove := make([]*websocket.Conn, 0)

	for client, clientUsername := range h.clients {
		if username == clientUsername {
			continue
		}
		err := client.WriteJSON(msg)
		if err != nil {
			// Vérifier si c'est une fermeture propre
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("ℹ️ Client déconnecté proprement :", err)
			} else {
				log.Println("❌ Erreur d'envoi du message de connexion :", err)
			}
			toRemove = append(toRemove, client) // Ajouter pour suppression
		}
	}

	// Supprimer directement les clients fermés
	for _, client := range toRemove {
		delete(h.clients, client)
	}
}
