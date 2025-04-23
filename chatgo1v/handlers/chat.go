package handlers

import (
	"log"
	"net/http"
	"simple-chat/db"
	"simple-chat/models"
	"time"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session, exists := GetSession(cookie.Value)
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := db.GetUserByID(session.UserID)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	mu.RLock()
	onlineCount := len(wsClients)
	mu.RUnlock()

	messages := db.GetMessages()

	data := struct {
		Username    string
		OnlineCount int
		Messages    []models.Message
	}{
		Username:    user.Username,
		OnlineCount: onlineCount,
		Messages:    messages,
	}

	err = tmpl.ExecuteTemplate(w, "chat.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ClearChatHandler(w http.ResponseWriter, r *http.Request) {
	db.ClearChat()

	mu.Lock()
	for client := range wsClients {
		if err := client.WriteJSON(map[string]string{
			"action": "clear_chat",
		}); err != nil {
			log.Printf("Error sending clear command: %v", err)
			client.Close()
			delete(wsClients, client)
		}
	}
	mu.Unlock()

	wsBroadcast <- models.Message{
		Username:  "System",
		Text:      "Чат был очищен",
		IsSystem:  true,
		Timestamp: time.Now(),
	}

	w.WriteHeader(http.StatusOK)
}
