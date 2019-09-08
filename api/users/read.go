package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "users", "remote_address": r.RemoteAddr, "path": "/api/users/{user}", "method": "GET"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		logger.WithField("user", vars["user"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	}
	logger.WithField("user", vars["user"]).Trace("Validated initial request")

	// Set searched user in logger
	logger = logger.WithField("user", vars["user"])

	// Check if user exists
	var user database.User
	db.Where("username = ?", vars["user"]).First(&user)
	if user.ID == 0 {
		logger.Trace("User does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}

	util.Responses.SuccessWithData(w, user)
	logger.Debug("Retrieved user from database")
}
