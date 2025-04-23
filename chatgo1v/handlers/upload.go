package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"simple-chat/db"
	"simple-chat/models"
	"strings"
	"time"
)

var imageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // 32MB max

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := "other"
	mimeType := mime.TypeByExtension(filepath.Ext(handler.Filename))
	if strings.HasPrefix(mimeType, "image/") {
		fileType = "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		fileType = "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		fileType = "audio"
	}

	ext := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("session_token")
	session, _ := GetSession(cookie.Value)
	user, _ := db.GetUserByID(session.UserID)

	msg := models.Message{
		Username:  user.Username,
		Text:      "/uploads/" + newFilename,
		Time:      time.Now().Format("15:04"),
		IsFile:    true,
		FileName:  handler.Filename,
		FileType:  fileType,
		Timestamp: time.Now(),
	}

	wsBroadcast <- msg
	w.WriteHeader(http.StatusOK)
}
