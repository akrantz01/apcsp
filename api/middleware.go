package main

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Handle authentication for all endpoints
func authMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow if authenticating
			if r.RequestURI == "/api/auth/login" || (r.RequestURI == "/api/users" && r.Method == "POST") {
				next.ServeHTTP(w, r)
				return
			}

			// Check if authorization header is present
			if r.Header.Get("Authorization") == "" {
				util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' must be present")
				return
			}

			// Validate JWT
			_, err := util.JWT.Validate(r.Header.Get("Authorization"), db)
			if err != nil {
				util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
