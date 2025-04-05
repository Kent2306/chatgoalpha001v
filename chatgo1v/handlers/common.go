package handlers

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"simple-chat/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Глобальные переменные пакета
var (
	// Шаблоны
	tmpl *template.Template

	// Сессии пользователей
	sessions = make(map[string]models.Session)
	mu       sync.RWMutex

	// Тема интерфейса
	theme = "light"

	// WebSocket
	wsClients   = make(map[*websocket.Conn]string)
	wsBroadcast = make(chan models.Message)
	wsUpgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// InitTemplates инициализирует HTML шаблоны
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

// CreateSession создает новую сессию пользователя
func CreateSession(userID int) models.Session {
	mu.Lock()
	defer mu.Unlock()

	token := uuid.New().String()
	session := models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sessions[token] = session
	return session
}

// GetSession возвращает сессию по токену
func GetSession(token string) (models.Session, bool) {
	mu.RLock()
	defer mu.RUnlock()
	session, exists := sessions[token]
	return session, exists
}

// ToggleTheme переключает между светлой и темной темой
func ToggleTheme() {
	mu.Lock()
	defer mu.Unlock()
	if theme == "light" {
		theme = "dark"
	} else {
		theme = "light"
	}
}

// GetTheme возвращает текущую тему
func GetTheme() string {
	mu.RLock()
	defer mu.RUnlock()
	return theme
}

// ToggleThemeHandler обрабатывает запрос на смену темы
func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	ToggleTheme()
	w.WriteHeader(http.StatusOK)
}
