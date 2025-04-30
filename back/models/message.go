package models

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	Content string          `json:"content"`
	Sender  *websocket.Conn `json:"-"`
}

type UserStatus struct {
	Type    string          `json:"type"`
	Content []string        `json:"content"`
	Sender  *websocket.Conn `json:"-"`
}

type PrivateMessage struct {
	Type    string          `json:"type"`
	Content []string        `json:"content"`
	Sender  *websocket.Conn `json:"-"`
}

type IsTyping struct {
	Type    string          `json:"type"`
	Content string          `json:"content"`
	Typing  bool            `json:"is_typing"`
	Sender  *websocket.Conn `json:"-"`
}
