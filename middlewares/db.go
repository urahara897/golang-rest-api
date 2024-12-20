package middlewares

import (
	"learn-golang/rest-api/db"
	"net/http"
)

func EnsureDB(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.InitDB()
		next(w, r)
	}
} 