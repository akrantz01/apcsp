package messages

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/akrantz01/apcsp/api/websockets"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func create(w http.ResponseWriter, r *http.Request, hub *websockets.Hub, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "messages", "remote_address": r.RemoteAddr, "path": "/api/chats/{chat}/messages", "method": "POST"})

	// Validate initial request on headers, path parameters, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		logger.WithField("chat", vars["chat"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		logger.WithFields(logrus.Fields{"chat": vars["chat"], "content_type": r.Header.Get("Content-Type")}).Trace("Invalid content type")
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		logger.WithField("chat", vars["chat"]).Trace("No request body given")
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}
	logger.WithField("chat", vars["chat"]).Trace("Validated initial request")

	// Add chat id to logger
	logger = logger.WithField("chat", vars["chat"])

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)
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

	// Ensure user is in chat
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

	// Validate JSON body
	var body struct {
		Type     string `json:"type"`
		Message  string `json:"message"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Type == "" {
		logger.WithField("type", body.Type).Trace("Field type not given")
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' is required")
		return
	} else if body.Type != "message" && body.Type != "image" && body.Type != "file" {
		logger.WithField("type", body.Type).Trace("Invalid type for 'type' field")
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be one of 'message', 'image', or 'file'")
		return
	} else if body.Type == "message" && body.Message == "" {
		logger.WithFields(logrus.Fields{"type": body.Type, "message": body.Message}).Trace("Message field must be present when type is 'message'")
		util.Responses.Error(w, http.StatusBadRequest, "field 'message' must be present")
		return
	} else if body.Type == "image" && body.Filename != "" {
		logger.WithFields(logrus.Fields{"type": body.Type, "filename": body.Filename}).Trace("Filename field not be present when type is 'filename'")
		util.Responses.Error(w, http.StatusBadRequest, "field 'filename' should be empty or nonexistent")
		return
	} else if body.Type == "file" && body.Filename == "" {
		logger.WithFields(logrus.Fields{"type": body.Type, "filename": body.Filename}).Trace("Filename filed must be present when type is 'filename'")
		util.Responses.Error(w, http.StatusBadRequest, "field 'filename' must be present")
		return
	}

	// Add chat id to logger
	logger = logger.WithField("uuid", chat.UUID)

	// Normal message
	if body.Type == "message" {
		// Get sender data by id
		var sender database.User
		db.Where("id = ?", uid).First(&sender)
		logger.Trace("Retrieved sender info from database")

		// Save message
		message := database.Message{
			ChatId:    chat.ID,
			SenderId:  uid,
			Type:      0,
			Message:   body.Message,
			Timestamp: time.Now().UnixNano(),
		}
		db.NewRecord(message)
		db.Create(&message)
		logger.Trace("Add message to database")

		// Associate with chat
		db.Model(&chat).Association("Messages").Append(&message)
		logger.Trace("Associate message with chat")

		// Push the message over websockets
		message.Sender = sender
		for _, user := range chat.Users {
			// Ignore sending user
			if user.ID == uid {
				continue
			}

			// Send message
			hub.PushMessage(user.Username, message, vars["chat"])
		}

		util.Responses.Success(w)
		logger.WithFields(logrus.Fields{"message": message.ID, "sender": message.SenderId}).Debug("Sent given message to chat")
		return
	}

	// Remove file name if image
	if body.Type == "image" {
		body.Filename = ""
		logger.WithField("type", body.Type).Trace("Removed filename for image message")
	}

	// Create file upload link
	id := uuid.NewV4().String()
	file := database.File{
		Path:     "./uploaded/" + id,
		Filename: body.Filename,
		UUID:     id,
		Used:     false,
		ChatId:   chat.ID,
	}
	db.NewRecord(file)
	db.Create(&file)
	logger.WithField("file", file.UUID).Trace("Created file link for message")

	// Create message database entry
	message := database.Message{
		ChatId:    chat.ID,
		SenderId:  uid,
		Type:      1,
		Message:   body.Message,
		FileId:    file.ID,
		Timestamp: time.Now().UnixNano(),
	}
	if body.Type == "file" {
		message.Message = ""
		message.Type = 2
		logger.WithField("type", body.Type).Trace("Change file type and empty message for file")
	}
	logger.WithField("file", file.UUID).Trace("Created message with file id")

	// Save to database
	db.NewRecord(message)
	db.Create(&message)
	logger.Trace("Saved message to database")

	// Associate with chat
	db.Model(&chat).Association("Messages").Append(&message)
	logger.Trace("Associate message with chat")

	util.Responses.SuccessWithData(w, map[string]string{"url": viper.GetString("http.domain") + "/api/files/" + file.UUID})
	logger.WithFields(logrus.Fields{"message": message.ID, "sender": message.SenderId, "file": file.UUID}).Debug("Created message with file upload link attached")
}
