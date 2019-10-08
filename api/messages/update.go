package messages

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "messages", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}/messages/{message}", "method": "PUT"})

	// Validate initial request on path parameters, headers, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for chat path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if _, ok := vars["message"]; !ok {
		logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]}).Trace("Invalid value for message path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'message' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	}
	logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]}).Trace("Validated initial request on path parameters")

	// Add chat id and message index to logger
	logger = logger.WithFields(logrus.Fields{"chat": vars["chat"], "message": vars["message"]})

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Where("uuid = ?", vars["chat"]).First(&chat)
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

	// Get message
	message := chat.Messages[index]

	// Parse JSON body
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	}

	// Modify message if passed
	if body.Message != "" {
		message.Message = body.Message
		message.Timestamp = time.Now().UnixNano()
		logger.Trace("Set new message for chat")
	}

	// Save changes
	db.Save(&message)
	logger.Trace("Saved updates to chat")

	util.Responses.Success(w)
	logger.Debug("Updated message in chat with specified data")
}
