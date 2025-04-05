package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	IsFile   bool   `json:"is_file"`
	FileName string `json:"file_name,omitempty"`
}
