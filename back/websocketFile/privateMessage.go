package websocketFile

import (
	"database/sql"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"real-time-forum/models"
)

func (h *Hub) BroadcastPrivateMessage(db *sql.DB, message, sender, receiver string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sender, err := getUsername(db, sender)
	if err != nil {
		log.Println(err)
		return
	}

	receiver, err = getUsername(db, receiver)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("BroadcastPrivateMessage called with message: %s, sender: %s, receiver: %s", message, sender, receiver)

	// Trouver les connexions des utilisateurs expéditeur et destinataire
	var senderConn, receiverConn *websocket.Conn

	for conn, username := range h.clients {
		if username == sender {
			senderConn = conn
			log.Printf("Sender connection found for user: %s", sender)
		} else if username == receiver {
			receiverConn = conn
			log.Printf("Receiver connection found for user: %s", receiver)
		}
	}

	// Vérifier si le destinataire est connecté
	if receiverConn == nil {
		log.Printf("Receiver (%s) not connected", receiver)
		return
	}

	// Créer le message privé
	privateMessage := models.PrivateMessage{
		Type:    "private",
		Content: []string{message, sender, time.Now().Format(time.RFC3339)},
		Sender:  senderConn,
	}
	log.Printf("Private message created: %+v", privateMessage)

	// Envoyer le message au destinataire
	if err := receiverConn.WriteJSON(privateMessage); err != nil {
		log.Printf("Error sending private message to %s: %v", receiver, err)
		receiverConn.Close()
		delete(h.clients, receiverConn)
	} else {
		log.Printf("Message sent successfully to %s", receiver)
	}
}

func getUsername(db *sql.DB, userID string) (string, error) {
	var username string
	query := `SELECT USERNAME FROM USER WHERE ID = ?`
	err := db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
