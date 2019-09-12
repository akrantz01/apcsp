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

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "messages", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}/messages"})

	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for path parameters")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}
	logger.WithField("chat", vars["chat"]).Trace("Validated initial request on path parameters")

	// Add chat id to logger
	logger = logger.WithField("chat", vars["chat"])

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Preload("Messages.Sender").Preload("Messages.File").Where("uuid = ?", vars["chat"]).First(&chat)
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

	page := int64(0)
	perPage := int64(100)
	if r.URL.Query().Get("page") != "" {
		page, err = strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			logger.WithError(err).WithField("page", r.URL.Query().Get("page")).Trace("Invalid page query parameter value")
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'page' must be an integer")
			return
		}
		logger.WithField("page", page).Trace("Set page to specified value in query parameter")
	}
	if r.URL.Query().Get("per_page") != "" {
		perPage, err = strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
		if err != nil {
			logger.WithError(err).WithField("per_page", r.URL.Query().Get("per_page")).Trace("Invalid per_page query parameter value")
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'per_page' must be an integer")
			return
		}
		logger.WithField("per_page", perPage).Trace("Set per_page to specified value in query parameter")
	}

	// Ensure length is not greater than end
	endIndex := perPage + (page * perPage)
	if endIndex > int64(len(chat.Messages)) {
		logger.Trace("Specified page out of bounds, defaulting to length of messages")
		endIndex = int64(len(chat.Messages))
	}
	logger.WithField("end", endIndex).Trace("Set end index for messages")

	// Isolate page of messages
	messages := chat.Messages[(page * perPage):endIndex]
	logger.Trace("Isolated page of messages based on page and number per page")

	// Assemble response map
	data := map[string]interface{}{
		"page":     page,
		"perPage":  perPage,
		"messages": messages,
	}
	logger.Trace("Created response map")

	util.Responses.SuccessWithData(w, data)
	logger.WithFields(logrus.Fields{"page": page, "per_page": perPage}).Debug("Got list of messages for specified chat")
}
