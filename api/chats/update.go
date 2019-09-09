package chats

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "chats", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}", "method": "PUT"})

	// Validate initial request on path parameters, headers, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		logger.WithFields(logrus.Fields{"chat": vars["chat"], "content_type": r.Header.Get("Content-Type")}).Trace("Invalid value for content type")
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		logger.WithField("chat", vars["chat"]).Trace("No request body given")
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	}
	logger.WithField("chat", vars["chat"]).Trace("Validated initial request")

	// Set chat id in logger
	logger = logger.WithField("chat", vars["chat"])

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)
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

	// Parse JSON body
	var body struct {
		Name string `json:"name"`
		Mode string `json:"mode"`
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Mode != "" && body.User == "" {
		logger.WithFields(logrus.Fields{"mode": body.Mode, "user": body.User}).Trace("Field user and mode must be passed together")
		util.Responses.Error(w, http.StatusBadRequest, "field 'user' must be passed when field 'mode' is present")
		return
	}

	// Modify name if passed
	if body.Name != "" {
		chat.DisplayName = body.Name
		logger.Trace("Set new display name for chat")
	}

	// Modify users associated with chat
	if body.Mode != "" {
		// Ensure user exists
		var user database.User
		db.Preload("Chats").Where("username = ?", body.User).First(&user)
		if user.ID == 0 {
			logger.WithField("user", body.User).Trace("Specified user to add/remove does not exist")
			util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
			return
		}
		logger.WithField("user", body.User).Trace("Retrieved user for addition/removal")

		switch body.Mode {
		// Add user to chat
		case "add":
			// Ensure not already in chat
			for _, c := range user.Chats {
				if c.ID == chat.ID {
					logger.WithField("user", body.User).Trace("Specified user already added")
					util.Responses.Error(w, http.StatusConflict, "specified user is already in chat")
					return
				}
			}

			// Add to chat
			db.Model(&chat).Association("Users").Append(&user)
			db.Model(&user).Association("Chats").Append(&chat)
			logger.WithField("user", body.User).Trace("Associated user with chat and chat with user")

		// Remove user from chat
		case "delete":
			// Ensure user is in chat
			valid := false
			for _, c := range user.Chats {
				if c.ID == chat.ID {
					valid = true
					break
				}
			}
			if !valid {
				logger.WithField("user", body.User).Trace("Specified user not in chat")
				util.Responses.Error(w, http.StatusBadRequest, "specified user is not in chat")
				return
			}

			db.Model(&chat).Association("Users").Delete(&user)
			db.Model(&user).Association("Chats").Delete(&chat)
			logger.WithField("user", body.User).Trace("Removed associated between user and chat")

		default:
			logger.WithField("mode", body.Mode).Trace("Invalid mode, must be add/delete")
			util.Responses.Error(w, http.StatusBadRequest, "field 'mode' must be one of 'add', 'delete'")
			return
		}
	}

	// Save changes
	db.Save(&chat)
	logger.Trace("Saved updates to chat")

	util.Responses.Success(w)
	logger.Debug("Updated chat with specified data")
}
