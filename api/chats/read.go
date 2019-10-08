package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func read(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "chats", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}", "method": "GET"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}
	logger.WithField("chat", vars["chat"]).Trace("Validated initial request")

	// Set chat id in logger
	logger = logger.WithField("chat", vars["chat"])

	// Check if chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Preload("Messages.Sender").Preload("Messages.File").Where("uuid = ?", vars["chat"]).First(&chat)
	if chat.ID == 0 {
		logger.Trace("Chat does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified chat does not exist")
		return
	}
	logger.Trace("Retrieved chat from database")

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

	// Check if requesting user is part of chat
	for _, user := range chat.Users {
		if id == user.ID {
			util.Responses.SuccessWithData(w, chat)
			logger.Debug("Retrieve chat from database")
			return
		}
	}

	util.Responses.Error(w, http.StatusForbidden, "user is not part of specified chat")
	logger.WithField("uid", id).Trace("User associated with token not in chat")
}
