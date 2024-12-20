package registration

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
	
	switch r.Method {
	case http.MethodPost:
		handleRegister(w, r)
	case http.MethodDelete:
		handleUnregister(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
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

	if err := event.Register(userID); err != nil {
		http.Error(w, "Could not register for event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully registered for event"})
}

func handleUnregister(w http.ResponseWriter, r *http.Request) {
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

	event := models.Event{ID: eventID}
	if err := event.CancelRegistration(userID); err != nil {
		http.Error(w, "Could not cancel registration", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully unregistered from event"})
} 