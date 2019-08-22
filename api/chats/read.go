package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func read(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}

	// Check if chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Preload("Messages.Sender").Where("uuid = ?", vars["chat"]).First(&chat)
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

	// Ensure id is of correct type
	idStr := util.JWT.Claims(token)["sub"]
	if _, ok := idStr.(string); !ok {
		util.Responses.Error(w, http.StatusBadRequest, "invalid type for 'subject' in token")
		return
	}

	// Ensure id is a number
	id, err := strconv.ParseUint(idStr.(string), 10, 32)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "invalid type for 'subject' in token")
	}

	// Check if requesting user is part of chat
	for _, user := range chat.Users {
		if uint(id) == user.ID {
			util.Responses.SuccessWithData(w, chat)
			return
		}
	}

	util.Responses.Error(w, http.StatusForbidden, "user is not part of specified chat")
}
