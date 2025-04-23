package main

import (
	"log"
	"net/http"
	"os"
	"simple-chat/handlers"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting server...")

	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("Template initialization failed: %v", err)
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	go handlers.HandleMessages()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/chat", handlers.AuthMiddleware(handlers.ChatHandler))
	r.HandleFunc("/ws", handlers.AuthMiddleware(handlers.HandleConnections))
	r.HandleFunc("/upload", handlers.AuthMiddleware(handlers.UploadHandler)).Methods("POST")
	r.HandleFunc("/clear-chat", handlers.AuthMiddleware(handlers.ClearChatHandler)).Methods("POST")
	r.HandleFunc("/online-users", handlers.AuthMiddleware(handlers.OnlineUsersHandler)).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
