package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func deleteMethod(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "chats", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}", "method": "DELETE"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}
	logger.WithField("chat", vars["chat"]).Trace("Validated initial request on path parameters")

	// Add chat id to logger
	logger = logger.WithField("chat", vars["chat"])

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)
	if chat.ID == 0 {
		logger.Trace("Specified chat does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified chat does not exist")
		return
	}
	logger.Trace("Retrieved chat information from database")

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		logger.WithError(err).Error("Unable to get unvalidated token")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}
	logger.Trace("Got unvalidated token parts")

	// Get user id from token
	id, err := util.JWT.UserId(token)
	if err != nil {
		logger.WithError(err).Trace("Failed to get user id from token")
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.WithField("uid", id).Trace("Got user id from token")

	// Get user from database
	var user database.User
	db.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
		return
	}
	logger.WithField("uid", id).Trace("Retrieved user data from database")

	// Ensure user is in chat
	valid := false
	for _, u := range chat.Users {
		if u.ID == id {
			valid = true
			break
		}
	}
	if !valid {
		logger.WithField("uid", id).Trace("User associated with token not in chat")
		util.Responses.Error(w, http.StatusForbidden, "specified user is not in chat")
		return
	}
	logger.WithField("uid", id).Trace("Confirmed requesting user in chat")

	// Delete chat and messages
	db.Delete(database.Message{}, "chat_id = ?", chat.ID)
	db.Delete(&chat)

	util.Responses.Success(w)
	logger.Debug("Deleted chat and all messages")
}
