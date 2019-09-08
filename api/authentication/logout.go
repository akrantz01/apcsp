package authentication

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Revoke authentication tokens
func Logout(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_addr": r.RemoteAddr, "path": "/api/auth/logout", "method": "GET"})

		// Get token w/o validation
		token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
		if err != nil {
			logger.WithError(err).Trace("Failed to decode token without validation")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
			return
		}
		logger.Trace("Got unvalidated token")

		// Delete token row from database
		var storedToken database.Token
		db.Where("id = ?", token.Header["kid"]).First(&storedToken)
		db.Delete(&storedToken)
		logger.Trace("Deleted given token")

		util.Responses.Success(w)
		logger.Debug("Revoked authentication token")
	}
}
