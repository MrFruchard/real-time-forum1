package websocketFile

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"real-time-forum/models"
	"sync"
)

// Hub gère toutes les connexions WebSocket actives et la diffusion des messages.
type Hub struct {
	clients    map[*websocket.Conn]string // Stocke les connexions WebSocket actives
	broadcast  chan interface{}           // Canal pour diffuser les messages à tous les clients
	register   chan *websocket.Conn       // Canal pour enregistrer une connexion
	unregister chan *websocket.Conn       // Canal pour supprimer une connexion
	mu         sync.Mutex                 // Mutex pour éviter les conflits d'accès
	DB         *sql.DB                    // Connexion à la base de données pour vérifier la session
}

// NewHub crée un nouveau hub WebSocket.
func NewHub(db *sql.DB) *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]string),
		broadcast:  make(chan interface{}),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		DB:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// On enregistre le client s'il n'est pas déjà présent
			if _, exists := h.clients[client]; !exists {
				h.clients[client] = "UNKNOWN_SESSION" // À remplacer par une vraie session si nécessaire
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			// On supprime le client s'il est présent
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.mu.Unlock()
				client.Close()
			} else {
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			// Conversion sécurisée du message
			msg, ok := message.(models.Message)
			if !ok {
				// Message de type inattendu, on l'ignore
				continue
			}

			// Récupérer une copie de la liste des clients à l'extérieur du verrou
			h.mu.Lock()
			clients := make([]*websocket.Conn, 0, len(h.clients))
			for client := range h.clients {
				// On n'envoie pas le message à l'expéditeur
				if client == msg.Sender {
					continue
				}
				clients = append(clients, client)
			}
			h.mu.Unlock()

			// Préparer le message à envoyer
			msgToSend := struct {
				Content string `json:"content"`
			}{
				Content: msg.Content,
			}

			// Diffuser le message en dehors de la section critique
			for _, client := range clients {
				if err := client.WriteJSON(msgToSend); err != nil {
					// En cas d'erreur, fermer la connexion et la retirer
					client.Close()
					h.mu.Lock()
					delete(h.clients, client)
					h.mu.Unlock()
				}
			}
		}
	}
}
