package handlers

import (
	"net/http"
	"simple-chat/db"
	"simple-chat/models"
	"time"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	session, exists := GetSession(cookie.Value)
	if !exists {
		return
	}

	user, err := db.GetUserByID(session.UserID)
	if err != nil {
		return
	}

	mu.Lock()
	wsClients[ws] = user.Username
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(wsClients, ws)
		mu.Unlock()
	}()

	for {
		var msg models.Message
		if err := ws.ReadJSON(&msg); err != nil {
			break
		}

		msg.Username = user.Username
		msg.Time = time.Now().Format("15:04")
		msg.Timestamp = time.Now()

		db.AddMessage(msg) // Теперь функция определена
		wsBroadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-wsBroadcast
		mu.RLock()
		for client := range wsClients {
			if err := client.WriteJSON(msg); err != nil {
				mu.RUnlock()
				mu.Lock()
				client.Close()
				delete(wsClients, client)
				mu.Unlock()
				mu.RLock()
			}
		}
		mu.RUnlock()
	}
}
