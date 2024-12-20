package events

import (
	"encoding/json"
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/models"
	"learn-golang/rest-api/utils"
	"net/http"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	db.InitDB()
	
	switch {
	case strings.Contains(r.URL.Path, "/events/"):
		handleSingleEvent(w, r)
	case r.Method == http.MethodGet:
		events, err := models.GetAllEvents()
		if err != nil {
			http.Error(w, "Could not fetch events", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(events)
		
	case r.Method == http.MethodPost:
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
	case r.Method == http.MethodPut:
		handleUpdateEvent(w, r)
	case r.Method == http.MethodDelete:
		handleDeleteEvent(w, r)
	}
}

func handleSingleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract event ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	
	eventID := parts[2]
	id, err := strconv.ParseInt(eventID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	
	event, err := models.GetEventByID(id)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(event)
}

func handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	eventID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	
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
	
	event, err := models.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	
	if *event.UserID != userID {
		http.Error(w, "Not authorized to update event", http.StatusUnauthorized)
		return
	}
	
	var updatedEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	
	updatedEvent.ID = eventID
	if err := updatedEvent.Update(); err != nil {
		http.Error(w, "Could not update event", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"message": "Event updated successfully"})
}

func handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	eventID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	
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
	
	event, err := models.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	
	if *event.UserID != userID {
		http.Error(w, "Not authorized to delete event", http.StatusUnauthorized)
		return
	}
	
	if err := event.Delete(); err != nil {
		http.Error(w, "Could not delete event", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"message": "Event deleted successfully"})
} 