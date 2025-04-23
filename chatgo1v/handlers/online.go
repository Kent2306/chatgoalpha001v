package handlers

import (
	"encoding/json"
	"net/http"
)

func OnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	users := make([]string, 0, len(wsClients))
	for _, username := range wsClients {
		users = append(users, username)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
