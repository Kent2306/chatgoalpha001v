package handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"simple-chat/models"
	"sync"
	"time"
)

var (
	tmpl *template.Template
	mu   sync.RWMutex

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsClients   = make(map[*websocket.Conn]string)
	wsBroadcast = make(chan models.Message)
	sessions    = make(map[string]models.Session)
)

func InitTemplates() error {
	var err error
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"firstChar": func(s string) string {
			if len(s) > 0 {
				return string(s[0])
			}
			return ""
		},
	}).ParseGlob("templates/*.html"))
	return err
}

func CreateSession(userID int) models.Session {
	mu.Lock()
	defer mu.Unlock()

	session := models.Session{
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sessions[session.Token] = session
	return session
}

func GetSession(token string) (models.Session, bool) {
	mu.RLock()
	defer mu.RUnlock()
	session, exists := sessions[token]
	return session, exists
}
