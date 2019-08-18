package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters and headers
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
		return
	} else if r.Header.Get("Authorization") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' must be present")
		return
	}

	// Ensure JWT is valid
	_, err := util.JWT.Validate(r.Header.Get("Authorization"), db)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
		return
	}

	// Check if user exists
	var user database.User
	db.Where("username = ?", vars["user"]).First(&user)
	if user.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}

	util.Responses.SuccessWithData(w, user)
}
