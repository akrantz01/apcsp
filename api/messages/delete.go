package messages

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func deleteMethod(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "messages", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}/messages/{message}", "method": "DELETE"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for chat path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if _, ok := vars["message"]; !ok {
		logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]}).Trace("Invalid value for message path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'message' must be present")
		return
	}
	logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]}).Trace("Validated initial request on path parameters")

	// Add chat id and message index to logger
	logger = logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]})

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Where("uuid = ?", vars["chat"]).First(&chat)
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
	uid, err := util.JWT.UserId(token)
	if err != nil {
		logger.WithError(err).Trace("Failed to get user id from token")
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.WithField("uid", uid).Trace("Got user id from token")

	// Check if requesting user is part of chat
	valid := false
	for _, user := range chat.Users {
		if uid == user.ID {
			valid = true
			break
		}
	}
	if !valid {
		logger.WithField("uid", uid).Trace("User associated with token not in chat")
		util.Responses.Error(w, http.StatusForbidden, "user is not part of specified chat")
		return
	}
	logger.WithField("uid", uid).Trace("Confirmed requesting user in chat")

	// Convert path parameter to integer
	index, err := strconv.ParseInt(vars["message"], 10, 64)
	if err != nil {
		logger.Trace("Invalid value for message index")
		util.Responses.Error(w, http.StatusBadRequest, "message index must be an integer")
		return
	}
	logger.Trace("Converted message index to integer")

	// Check integer bounds
	if index > int64(len(chat.Messages)-1) {
		logger.WithField("total_messages", len(chat.Messages)-1).Trace("Message index greater than total messages")
		util.Responses.Error(w, http.StatusBadRequest, "message index is out of bounds")
		return
	}
	logger.Trace("Validated message index not out of bounds")

	// Delete specified message
	db.Delete(&chat.Messages[index])

	util.Responses.Success(w)
	logger.Debug("Delete message at specified index")
}
