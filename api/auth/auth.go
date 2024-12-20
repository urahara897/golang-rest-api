package auth

import (
	"encoding/json"
	"learn-golang/rest-api/db"
	"learn-golang/rest-api/models"
	"learn-golang/rest-api/utils"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	db.InitDB()
	
	path := r.URL.Path
	
	switch {
	case strings.HasSuffix(path, "/signup"):
		handleSignup(w, r)
	case strings.HasSuffix(path, "/login"):
		handleLogin(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := user.Save(); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	id, err := user.ValidateCredentials()
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.Email, id)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
} 