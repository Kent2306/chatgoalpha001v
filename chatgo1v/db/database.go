package db

import (
	"fmt"
	"simple-chat/models"
	"sync"
	"time"
)

var (
	users     = make(map[string]models.User)
	messages  []models.Message
	dbMu      sync.RWMutex
	nextID    = 1
	nextMsgID = 1
)

func CreateUser(user models.User) error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if _, exists := users[user.Username]; exists {
		return fmt.Errorf("username already exists")
	}

	user.ID = nextID
	nextID++
	users[user.Username] = user
	return nil
}

func GetUser(username string) (models.User, bool) {
	dbMu.RLock()
	defer dbMu.RUnlock()
	user, exists := users[username]
	return user, exists
}

func GetUserByID(id int) (models.User, error) {
	dbMu.RLock()
	defer dbMu.RUnlock()
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, fmt.Errorf("user not found")
}

func AddMessage(msg models.Message) models.Message {
	dbMu.Lock()
	defer dbMu.Unlock()

	msg.ID = nextMsgID
	nextMsgID++
	msg.Timestamp = time.Now()
	messages = append(messages, msg)
	return msg
}

func GetMessages() []models.Message {
	dbMu.RLock()
	defer dbMu.RUnlock()
	messagesCopy := make([]models.Message, len(messages))
	copy(messagesCopy, messages)
	return messagesCopy
}

func ClearChat() {
	dbMu.Lock()
	defer dbMu.Unlock()
	messages = []models.Message{}
	nextMsgID = 1
}
