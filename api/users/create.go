package users

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on Content-Type header and body
	if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}

	// Validate JSON body
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Username == "" || body.Password == "" || body.Email == "" || body.Name == "" {
		util.Responses.Error(w, http.StatusBadRequest, "fields 'username', 'password', 'email', and 'name' are required")
		return
	} else if len(body.Password) != 64 {
		util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
		return
	} else if len(body.Email) < 5 || len(body.Email) > 254 {
		util.Responses.Error(w, http.StatusBadRequest, "field 'email' must be of length between 5 and 254 characters")
		return
	}

	// Check if username is taken
	var user database.User
	db.Where("username = ?", body.Username).First(&user)
	if user.ID != 0 {
		util.Responses.Error(w, http.StatusConflict, "username is already taken")
		return
	}

	// Check if email is taken
	db.Where("email = ?", body.Email).First(&user)
	if user.ID != 0 {
		util.Responses.Error(w, http.StatusConflict, "email is already taken")
		return
	}

	// Hash the password for security
	hash, err := passlib.Hash(body.Password)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	// Assemble user info from body
	user.Name = body.Name
	user.Email = body.Email
	user.Username = body.Username
	user.Password = hash

	// Save to database
	db.NewRecord(user)
	db.Save(&user)

	util.Responses.Success(w)
}
