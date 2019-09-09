package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func deleteMethod(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "users", "remote_address": r.RemoteAddr, "path": "/api/users/{user}", "method": "DELETE"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		logger.WithField("user", vars["user"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	}
	logger.WithField("user", vars["user"]).Trace("Validated initial request")

	// Add username to logger
	logger = logger.WithField("user", vars["user"])

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		logger.WithError(err).Error("Unable to get unvalidated token")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}
	logger.Trace("Got unvalidated token parts")

	// Get user from database
	var user database.User
	db.Where("username = ?", vars["user"]).First(&user)
	if user.ID == 0 {
		logger.Trace("User specified does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}
	logger.Trace("Retrieved user information from database")

	// Ensure user from token is user being modified
	if sameUser, err := util.JWT.CheckUser(token, user, db); err != nil {
		logger.WithError(err).Trace("User specified in token not found")
		util.Responses.Error(w, http.StatusBadRequest, "user associated with token not found")
		return
	} else if !sameUser {
		logger.Trace("User in token and user specified do not match")
		util.Responses.Error(w, http.StatusForbidden, "not allowed to modify other users")
		return
	}

	// Delete the user and all associated tokens
	db.Delete(database.Token{}, "user_id = ?", user.ID)
	db.Delete(&user)
	logger.Trace("Revoked all tokens and deleted user")

	util.Responses.Success(w)
	logger.Debug("Deleted all user information")
}
