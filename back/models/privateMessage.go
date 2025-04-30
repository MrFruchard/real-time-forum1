package models

type PrivateMessageReceived struct {
	Sender   string
	Message  string `json:"message"`
	Receiver string `json:"receiver"`
}

type PrivateMessageSend struct {
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	ReceiverID     string `json:"receiver_id"`
	Message        string `json:"message"`
	MessageID      string `json:"message_id"`
}

type SendAllPrivateMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Date    string `json:"date"`
}
