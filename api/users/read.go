package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["user"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'user' must be present")
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
