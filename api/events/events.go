package events

import (
	"encoding/json"
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/models"
	"learn-golang/rest-api/utils"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	db.InitDB()
	
	switch r.Method {
	case http.MethodGet:
		events, err := models.GetAllEvents()
		if err != nil {
			http.Error(w, "Could not fetch events", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(events)
		
	case http.MethodPost:
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		
		userID, err := utils.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var event models.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		
		event.UserID = &userID
		if err := event.Save(); err != nil {
			http.Error(w, "Could not create event", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Event created",
			"event":   event,
		})
	}
} 