package models

type GetNotification struct {
	Id        int    `json:"id"`
	Sender    string `json:"sender"`
	Type      string `json:"type"`
	RelatedId int    `json:"related_id"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}
