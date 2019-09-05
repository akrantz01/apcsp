package messages

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	}

	// Check that chat exists
	var chat database.Chat
	db.Preload("Users").Preload("Messages").Preload("Messages.Sender").Preload("Messages.File").Where("uuid = ?", vars["chat"]).First(&chat)
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

	page := int64(0)
	perPage := int64(100)
	if r.URL.Query().Get("page") != "" {
		page, err = strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'page' must be an integer")
			return
		}
	}
	if r.URL.Query().Get("per_page") != "" {
		perPage, err = strconv.ParseInt(r.URL.Query().Get("per_page"), 10, 64)
		if err != nil {
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'per_page' must be an integer")
			return
		}
	}

	// Ensure length is not greater than end
	endIndex := perPage + (page * perPage)
	if endIndex > int64(len(chat.Messages)) {
		endIndex = int64(len(chat.Messages))
	}

	// Isolate page of messages
	messages := chat.Messages[(page * perPage):endIndex]

	// Assemble response map
	data := map[string]interface{}{
		"page":     page,
		"perPage":  perPage,
		"messages": messages,
	}

	util.Responses.SuccessWithData(w, data)
}
