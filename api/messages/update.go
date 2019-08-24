package messages

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters, headers, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if _, ok := vars["message"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'message' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "body must be present")
		return
	}

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Where("uuid = ?", vars["chat"]).First(&chat)
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
	uid, err := util.JWT.UserId(token)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if requesting user is part of chat
	valid := false
	for _, user := range chat.Users {
		if uid == user.ID {
			valid = true
			break
		}
	}
	if !valid {
		util.Responses.Error(w, http.StatusForbidden, "user is not part of specified chat")
		return
	}

	// Convert path parameter to integer
	index, err := strconv.ParseInt(vars["message"], 10, 64)
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "message index must be an integer")
		return
	}

	// Check integer bounds
	if index > int64(len(chat.Messages)-1) {
		util.Responses.Error(w, http.StatusBadRequest, "message index is out of bounds")
		return
	}

	// Get message
	message := chat.Messages[index]

	// Parse JSON body
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	}

	// Modify message if passed
	if body.Message != "" {
		message.Message = body.Message
		message.Timestamp = time.Now().UnixNano()
	}

	db.Save(&message)

	util.Responses.Success(w)
}
