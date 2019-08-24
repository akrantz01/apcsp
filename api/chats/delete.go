package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func deleteMethod(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)
	if chat.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified chat does not exist")
		return
	}

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}

	// Get user id from token
	id, err := util.JWT.UserId(token)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user from database
	var user database.User
	db.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}

	// Ensure user is in chat
	valid := false
	for _, u := range chat.Users {
		if u.ID == id {
			valid = true
			break
		}
	}
	if !valid {
		util.Responses.Error(w, http.StatusForbidden, "specified user is not in chat")
		return
	}

	// Delete chat and messages
	db.Delete(database.Message{}, "chat_id = ?", chat.ID)
	db.Delete(&chat)

	util.Responses.Success(w)
}
