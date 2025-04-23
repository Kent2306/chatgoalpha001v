package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Time      string    `json:"time"`
	IsFile    bool      `json:"is_file"`
	FileName  string    `json:"file_name,omitempty"`
	IsSystem  bool      `json:"is_system"`
	FileType  string    `json:"file_type,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
