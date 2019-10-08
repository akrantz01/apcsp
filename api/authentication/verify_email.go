package authentication

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func VerifyEmail(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_address": r.RemoteAddr, "path": "/api/auth/verify-email", "method": "GET"})

		// Validate request on method and query parameters
		if r.Method != http.MethodGet {
			logger.WithField("method", r.Method).Trace("Invalid request method")
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if len(r.URL.RawQuery) == 0 || r.URL.Query().Get("token") == "" {
			logger.WithField("token", r.URL.Query().Get("token")).Trace("Invalid query parameter")
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'token' is required")
			return
		}
		logger.Trace("Validated request")

		// Validate JWT
		token, err := util.JWT.Validate(r.URL.Query().Get("token"), database.TokenVerification, db)
		if err != nil {
			util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
			return
		}
		logger.Trace("Successfully validated authentication token")

		// Get user from token
		userId, _ := util.JWT.UserId(token)
		var user database.User
		db.Where("id = ?", userId).First(&user)
		if user.ID == 0 {
			logger.WithField("id", userId).Trace("User not found in database")
			util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
			return
		}

		// Add user to logger
		logger = logger.WithField("username", user.Username)

		// Set as verified
		user.Verified = true
		db.Save(&user)

		// Delete verification token
		db.Delete(database.Token{}, "id = ?", token.Header["kid"])

		util.Responses.Success(w)
		logger.Debug("Successfully verified user email")
	}
}
