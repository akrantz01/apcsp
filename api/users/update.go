package users

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters, headers and body
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	}

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}

	// Get user from database
	var user database.User
	db.Where("username = ?", vars["user"]).First(&user)
	if user.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}

	// Ensure user from token is user being modified
	if sameUser, err := util.JWT.CheckUser(token, user, db); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "user associated with token not found")
		return
	} else if !sameUser {
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
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	}

	// Modify name if passed
	if body.Name != "" {
		user.Name = body.Name
	}

	// Modify email if passed
	if body.Email != "" {
		// Validate length
		if len(body.Email) < 5 || len(body.Email) > 254 {
			util.Responses.Error(w, http.StatusBadRequest, "field 'email' must be of length between 5 and 254")
			return
		}
		user.Email = body.Email
	}

	// Modify password if passed
	if body.Password != "" {
		// Validate length
		if len(body.Password) != 64 {
			util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
			return
		}

		// Hash password
		hash, err := passlib.Hash(body.Password)
		if err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		user.Password = hash
	}

	// Save changes
	db.Save(&user)

	util.Responses.Success(w)
}
