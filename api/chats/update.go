package chats

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func update(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path parameters, headers, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
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
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)

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
	valid := false
	for _, user := range chat.Users {
		if uint(id) == user.ID {
			valid = true
			break
		}
	}
	if !valid {
		util.Responses.Error(w, http.StatusUnauthorized, "user is not part of specified chat")
		return
	}

	// Parse JSON body
	var body struct {
		Name string `json:"name"`
		Mode string `json:"mode"`
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Mode != "" && body.User == "" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'users' must be passed when field 'mode' is present")
		return
	}

	// Modify name if passed
	if body.Name != "" {
		chat.DisplayName = body.Name
	}

	// Modify users associated with chat
	if body.Mode != "" {
		// Ensure user exists
		var user database.User
		db.Preload("Chats").Where("username = ?", body.User).First(&user)
		if user.ID == 0 {
			util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
			return
		}

		switch body.Mode {
		// Add user to chat
		case "add":
			// Ensure not already in chat
			for _, c := range user.Chats {
				if c.ID == chat.ID {
					util.Responses.Error(w, http.StatusConflict, "specified user is already in chat")
					return
				}
			}

			// Add to chat
			db.Model(&chat).Association("Users").Append(&user)
			db.Model(&user).Association("Chats").Append(&chat)

		// Remove user from chat
		case "delete":
			// Ensure user is in chat
			valid := false
			for _, c := range user.Chats {
				if c.ID == chat.ID {
					valid = true
				}
			}
			if !valid {
				util.Responses.Error(w, http.StatusBadRequest, "specified user is not in chat")
				return
			}

			db.Model(&chat).Association("Users").Delete(&user)
			db.Model(&user).Association("Chats").Delete(&chat)

		default:
			util.Responses.Error(w, http.StatusBadRequest, "field 'mode' must be one of 'add', 'delete'")
			return
		}
	}

	// Save changes
	db.Save(&chat)

	util.Responses.Success(w)
}
