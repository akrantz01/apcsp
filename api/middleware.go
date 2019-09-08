package main

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = logrus.WithField("app", "middleware")

// Handle authentication for all endpoints
func authMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set remote address to X-Forwarded-For or CF-Connecting-IP if available
			if r.Header.Get("CF-Connecting-IP") != "" {
				r.RemoteAddr = r.Header.Get("CF-Connecting-IP")
				logger.WithField("remote-address", r.RemoteAddr).Trace("Setting CF-Connecting-IP to remote address")
			} else if r.Header.Get("X-Forwarded-For") != "" {
				r.RemoteAddr = r.Header.Get("X-Forwarded-For")
				logger.WithField("remote-address", r.RemoteAddr).Trace("Setting X-Forwarded-For to remote address")
			}

			// Allow if authenticating
			if r.RequestURI == "/api/auth/login" || (r.RequestURI == "/api/users" && r.Method == "POST") {
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

			// Validate JWT
			_, err := util.JWT.Validate(r.Header.Get("Authorization"), db)
			if err != nil {
				util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
				return
			}
			logger.Trace("Successfully validated authentication token")

			next.ServeHTTP(w, r)
		})
	}
}
