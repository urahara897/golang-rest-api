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
	case strings.HasPrefix(path, "/admin/signup"):
		handleAdminSignup(w, r)
	case strings.HasPrefix(path, "/admin/login"):
		handleAdminLogin(w, r)
	case strings.HasSuffix(path, "/users"):
		handleUsers(w, r)
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

func handleAdminSignup(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := admin.Save(); err != nil {
		http.Error(w, "Could not create admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Admin created successfully",
		"admin":   admin,
	})
}

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	id, err := admin.ValidateCredentials()
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(admin.Email, id)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
} 