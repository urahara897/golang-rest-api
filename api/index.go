package api

import (
	"encoding/json"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Your handler logic here
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from Go API",
	})
} 