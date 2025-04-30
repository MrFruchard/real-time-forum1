package models

import "time"

type Post struct {
	ID        int
	UserId    string
	Content   string `json:"content"`
	Likes     int
	Dislikes  int
	Comments  []Comment
	Category  string `json:"category"`
	Image     string
	CreatedAt time.Time
}
