package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"simple-chat/db"
	"simple-chat/models"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим multipart форму (макс. 10MB)
	r.ParseMultipartForm(10 << 20)

	// Получаем файл из запроса
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем уникальное имя файла
	ext := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Создаем новый файл на сервере
	dst, err := os.Create("uploads/" + newFilename)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о пользователе
	cookie, _ := r.Cookie("session_token")
	session, _ := GetSession(cookie.Value)
	user, _ := db.GetUserByID(session.UserID)

	// Создаем сообщение о загрузке файла
	msg := models.Message{
		Username: user.Username,
		Text:     fmt.Sprintf("/uploads/%s", newFilename),
		Time:     time.Now().Format("15:04"),
		IsFile:   true,
		FileName: handler.Filename,
	}

	// Отправляем сообщение в канал WebSocket
	wsBroadcast <- msg
	w.WriteHeader(http.StatusOK)
}
