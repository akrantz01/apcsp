package chats

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "chats", "remote_address": r.RemoteAddr, "path": "/api/chats", "method": "POST"})

	// Validate initial request on headers and body
	if r.Header.Get("Content-Type") != "application/json" {
		logger.WithField("content_type", r.Header.Get("Content-Type")).Trace("Invalid content type")
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		logger.Trace("No request body given")
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
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Name == "" || len(body.Users) == 0 || body.Message == "" {
		logger.WithFields(logrus.Fields{"name": body.Name, "users": body.Users, "message": body.Message}).Trace("Field name, users, or message not given")
		util.Responses.Error(w, http.StatusBadRequest, "fields 'name', 'users', and 'message' are required")
		return
	}

	// Add chat name to logger
	logger = logger.WithField("name", body.Name)

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		logger.WithError(err).Error("Unable to get unvalidated token")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}
	logger.Trace("Got unvalidated token parts")

	// Get requesting user from database
	var requestingUser database.User
	db.Where("id = ?", util.JWT.Claims(token)["sub"]).First(&requestingUser)
	logger.WithFields(logrus.Fields{"username": requestingUser.Username, "id": requestingUser.ID}).Trace("Retrieved requesting user from database")

	// Check all users specified exist
	var users []*database.User
	for _, name := range body.Users {
		// Ensure not requesting self
		if requestingUser.Username == name {
			logger.Trace("Requesting user attempted to add self to chat")
			util.Responses.Error(w, http.StatusBadRequest, "current user cannot be included in recipients")
			return
		}

		// Check user exists
		var checkUser database.User
		db.Where("username = ?", name).First(&checkUser)
		if checkUser.ID == 0 {
			logger.WithField("username", name).Trace("Cannot add user to chat, does not exist")
			util.Responses.Error(w, http.StatusBadRequest, "username '"+name+"' does not exist")
			return
		}
		users = append(users, &checkUser)
		logger.Trace("Added user data to chat users")
	}
	users = append(users, &requestingUser)
	logger.Trace("Automatically added requesting user to chat users")

	// Generate UUID for chat
	u, err := uuid.NewV4().MarshalText()
	if err != nil {
		logger.WithError(err).Trace("Failed to encode UUID")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to encode uuid to text")
		return
	}
	logger.WithField("uuid", u).Trace("Assign id to chat")

	// Add uuid to chat
	logger = logger.WithField("uuid", u)

	// Create chat
	chat := &database.Chat{
		DisplayName: body.Name,
		UUID:        string(u),
	}
	db.NewRecord(chat)
	db.Create(&chat)
	logger.WithFields(logrus.Fields{"name": chat.DisplayName, "users": body.Users}).Trace("Create chat with users and name")

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
	logger.Trace("Generate initial message for chat")

	// Associate with messages
	chat.Messages = append(chat.Messages, message)
	db.Save(&chat)
	logger.Trace("Add message to chat")

	// Add chat to each user and each user to chat
	for _, user := range users {
		db.Model(&user).Association("Chats").Append(chat)
		db.Model(&chat).Association("Users").Append(user)
	}
	logger.Trace("Add all users to chat and chat to users")

	util.Responses.Success(w)
	logger.WithFields(logrus.Fields{"name": chat.DisplayName, "users": body.Users}).Debug("Created chat with name, users, and initial message")
}
