package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func deleteMethod(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Get token w/o validation
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
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
	db.Where("username = ?", vars["chat"]).First(&user)
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

	// Delete the
	vars := mux.Vars(r)
	db.Delete(database.Chat{}, "id = ?", vars["chat"])



}
