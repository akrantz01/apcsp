package chats

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on headers and body
	if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}

	// Validate JSON body
	var body struct {
		Name    string   `json:"name"`
		Users   []string `json:"users"`
		Message string   `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Name == "" || len(body.Users) == 0 || body.Message == "" {
		util.Responses.Error(w, http.StatusBadRequest, "fields 'name', 'users', and 'message' are required")
		return
	}

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}

	// Get requesting user from database
	var requestingUser database.User
	db.Where("id = ?", util.JWT.Claims(token)["sub"]).First(&requestingUser)

	// Check all users specified exist
	var users []*database.User
	for _, name := range body.Users {
		// Ensure not requesting self
		if requestingUser.Username == name {
			util.Responses.Error(w, http.StatusBadRequest, "current user cannot be included in recipients")
			return
		}

		// Check user exists
		var checkUser database.User
		db.Where("username = ?", name).First(&checkUser)
		if checkUser.ID == 0 {
			util.Responses.Error(w, http.StatusBadRequest, "username '"+name+"' does not exist")
			return
		}
		users = append(users, &checkUser)
	}
	users = append(users, &requestingUser)

	// Create chat
	chat := &database.Chat{
		Name: body.Name,
	}
	db.NewRecord(chat)
	db.Create(&chat)

	// Add initial message to chat
	message := database.Message{
		ChatId:    chat.ID,
		Sender:    requestingUser,
		SenderId:  requestingUser.ID,
		Message:   body.Message,
		Timestamp: time.Now().UnixNano(),
	}
	db.NewRecord(message)
	db.Create(&message)

	// Associate with messages
	chat.Messages = append(chat.Messages, message)
	db.Save(&chat)

	// Add chat to each user and each user to chat
	for _, user := range users {
		db.Model(&user).Association("Chats").Append(chat)
		db.Model(&chat).Association("Users").Append(user)
	}

	util.Responses.Success(w)
}
