package handlers

import (
	"log"
	"net/http"
	"time"

	"simple-chat/db"
	"simple-chat/models"

	"github.com/gorilla/websocket"
)

// HandleConnections устанавливает WebSocket соединение и обрабатывает входящие сообщения
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Error getting cookie: %v", err)
		return
	}

	// Получаем сессию пользователя
	session, exists := GetSession(cookie.Value)
	if !exists {
		log.Printf("Session not found for token: %s", cookie.Value)
		return
	}

	// Получаем данные пользователя из БД
	user, err := db.GetUserByID(session.UserID)
	if err != nil {
		log.Printf("Error getting user by ID %d: %v", session.UserID, err)
		return
	}

	// Обновляем соединение до WebSocket
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer func() {
		ws.Close()
		delete(wsClients, ws)
	}()

	// Регистрируем нового клиента
	wsClients[ws] = user.Username
	log.Printf("New WebSocket connection from %s", user.Username)

	// Бесконечный цикл обработки сообщений
	for {
		var msg models.Message

		// Читаем JSON сообщение
		if err := ws.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Добавляем метаданные
		msg.Username = user.Username
		msg.Time = time.Now().Format("15:04")
		log.Printf("Received message from %s: %s", user.Username, msg.Text)

		// Отправляем в канал трансляции
		wsBroadcast <- msg
	}
}

// HandleMessages рассылает сообщения всем подключенным клиентам
func HandleMessages() {
	for {
		// Получаем сообщение из канала
		msg := <-wsBroadcast

		// Рассылаем всем подключенным клиентам
		for client := range wsClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket write error to %s: %v", wsClients[client], err)
				client.Close()
				delete(wsClients, client)
			} else {
				log.Printf("Sent message to %s", wsClients[client])
			}
		}
	}
}
