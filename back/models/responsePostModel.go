package models

type ResponsePost struct {
	Id        int    `json:"id"`
	User      string `json:"username"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	Comments  int    `json:"comments"`
	Category  string `json:"category"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	Liked     bool   `json:"liked"`
	Disliked  bool   `json:"disliked"`
}

type ResponseComment struct {
	Id        int    `json:"id"`
	User      string `json:"username"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	CreatedAt string `json:"created_at"`
	Liked     bool   `json:"liked"`
	Disliked  bool   `json:"disliked"`
}
