package main

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Handle authentication for all endpoints
func authMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logrus.WithField("app", "middleware")

			// Add remote address to logger
			logger = logrus.WithFields(logrus.Fields{"app": "middleware", "remote_address": r.RemoteAddr})

			// Allow if authenticating
			if r.RequestURI == "/api/auth/login" || (r.RequestURI == "/api/users" && r.Method == "POST") || r.RequestURI == "/api/ws" || strings.Index(r.RequestURI, "/api/auth/forgot-password") == 0 || strings.Index(r.RequestURI, "/api/auth/verify-email") == 0 || strings.Index(r.RequestURI, "/api/") == -1 {
				logger.WithField("uri", r.RequestURI).Trace("Unauthenticated route received")
				next.ServeHTTP(w, r)
				return
			}
			logger.Trace("Authentication for route required")

			// Check if authorization header is present
			if r.Header.Get("Authorization") == "" {
				util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' must be present")
				return
			}
			logger.Trace("Ensured authentication header is present")

			// Select token type
			tokenType := database.TokenAuthentication
			if strings.Index(r.RequestURI, "/api/auth/reset-password") == 0 {
				tokenType = database.TokenResetPassword
			}

			// Validate JWT
			_, err := util.JWT.Validate(r.Header.Get("Authorization"), tokenType, db)
			if err != nil {
				util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}
			logger.Trace("Successfully validated authentication token")

			next.ServeHTTP(w, r)
		})
	}
}
