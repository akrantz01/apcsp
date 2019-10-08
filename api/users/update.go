package users

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "users", "remote_address": r.RemoteAddr, "path": "/api/users/{user}", "method": "PUT"})

	// Validate initial request on path parameters, headers and body
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		logger.WithField("user", vars["user"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		logger.WithFields(logrus.Fields{"user": vars["user"], "content_type": r.Header.Get("Content-Type")}).Trace("Invalid value for content type header")
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		logger.WithField("user", vars["user"]).Trace("No request body given")
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	}
	logger.WithField("user", r.URL.Query().Get("username")).Trace("Validated initial request")

	// Set searched user in logger
	logger = logger.WithField("username", r.URL.Query().Get("username"))

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
		logger.Trace("User does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}
	logger.Trace("Retrieved user from database")

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

	// Parse JSON body
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	}
	logger.Trace("Validated JSON body")

	// Modify name if passed
	if body.Name != "" {
		user.Name = body.Name
		logger.Trace("Set new name for user")
	}

	// Modify email if passed
	if body.Email != "" {
		// Validate length
		if len(body.Email) < 5 || len(body.Email) > 254 {
			logger.WithField("email", len(body.Email)).Trace("Invalid email length")
			util.Responses.Error(w, http.StatusBadRequest, "field 'email' must be of length between 5 and 254")
			return
		}
		user.Email = body.Email
		logger.Trace("Set new email for user")
	}

	// Modify password if passed
	if body.Password != "" {
		// Validate length
		if len(body.Password) != 64 {
			logger.WithField("password", len(body.Password)).Trace("Invalid password length")
			util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
			return
		}

		// Hash password
		hash, err := passlib.Hash(body.Password)
		if err != nil {
			logger.WithError(err).Error("Failed to hash password")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		user.Password = hash
		logger.Trace("Set new password for user")
	}

	// Save changes
	db.Save(&user)
	logger.Trace("Saved new user information to database")

	util.Responses.Success(w)
	logger.Debug("Updated user with specified data")
}
