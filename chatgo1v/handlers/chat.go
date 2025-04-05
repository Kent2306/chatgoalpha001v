package handlers

import (
	"context"
	"net/http"
	"simple-chat/db"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	user, _ := db.GetUserByID(userID)

	data := struct {
		Username string
		Theme    string
	}{
		Username: user.Username,
		Theme:    GetTheme(),
	}

	if err := tmpl.ExecuteTemplate(w, "chat.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
