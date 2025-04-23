package handlers

import (
	"net/http"
	"simple-chat/db"
	"simple-chat/models"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Дополнительная структура для отслеживания первого входа
var (
	firstLogin     = make(map[int]bool)
	firstLoginLock sync.RWMutex
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, "Ошибка при загрузке страницы входа", http.StatusInternalServerError)
		}
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, exists := db.GetUser(username)
	if !exists {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	// Создаем сессию и отмечаем первый вход
	session := CreateSession(user.ID)

	firstLoginLock.Lock()
	firstLogin[user.ID] = true
	firstLoginLock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.ExpiresAt,
		Path:    "/",
	})

	http.Redirect(w, r, "/chat", http.StatusFound)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, "Ошибка при загрузке страницы регистрации", http.StatusInternalServerError)
		}
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(username) < 3 || len(username) > 20 {
		http.Error(w, "Имя пользователя должно быть от 3 до 20 символов", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Пароль должен содержать минимум 8 символов", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при создании учетной записи", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := db.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		mu.Lock()
		delete(sessions, cookie.Value)
		mu.Unlock()

		// Получаем ID пользователя из сессии
		if session, exists := GetSession(cookie.Value); exists {
			firstLoginLock.Lock()
			delete(firstLogin, session.UserID)
			firstLoginLock.Unlock()
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if _, exists := GetSession(cookie.Value); !exists {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// SendWelcomeMessage отправляет приветственное сообщение при первом входе
func SendWelcomeMessage(userID int) {
	firstLoginLock.RLock()
	defer firstLoginLock.RUnlock()

	if firstLogin[userID] {
		wsBroadcast <- models.Message{
			Username:  "System",
			Text:      "Добро пожаловать в SunSet! Помните: 1) Это анонимный чат 2) Администрация не несёт ответственности 3) Запрещены незаконные действия",
			IsSystem:  true,
			Timestamp: time.Now(),
		}

		firstLoginLock.Lock()
		firstLogin[userID] = false
		firstLoginLock.Unlock()
	}
}
