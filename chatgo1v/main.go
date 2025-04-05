package main

import (
	"log"
	"net/http"
	"os"
	"simple-chat/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Создаем папку для загрузок, если ее нет
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Инициализируем шаблоны
	if err := handlers.InitTemplates(); err != nil {
		log.Fatal("Error initializing templates:", err)
	}

	// Запускаем обработчик WebSocket сообщений в отдельной горутине
	go handlers.HandleMessages()

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Настраиваем обработку статических файлов
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Регистрируем маршруты
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/chat", handlers.AuthMiddleware(handlers.ChatHandler))
	r.HandleFunc("/ws", handlers.AuthMiddleware(handlers.HandleConnections))
	r.HandleFunc("/upload", handlers.AuthMiddleware(handlers.UploadHandler)).Methods("POST")
	r.HandleFunc("/toggle-theme", handlers.ToggleThemeHandler).Methods("POST")

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
