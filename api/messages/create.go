package messages

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on headers, path parameters, and body
	vars := mux.Vars(r)
	if _, ok := vars["chat"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'chat' must be present")
		return
	} else if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("uuid = ?", vars["chat"]).First(&chat)
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

	// Ensure user is in chat
	valid := false
	for _, chat := range chat.Users {
		if uid == chat.ID {
			valid = true
			break
		}
	}
	if !valid {
		util.Responses.Error(w, http.StatusForbidden, "user is not part of specified chat")
		return
	}

	// Validate JSON body
	var body struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Type == "" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' is required")
		return
	} else if body.Type != "message" && body.Type != "image" && body.Type != "file" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'type' must be one of 'message', 'image', or 'file'")
		return
	} else if body.Type == "message" && body.Message == "" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'message' must be present")
		return
	} else if body.Type == "image" && body.Filename != "" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'filename' should be empty or nonexistent")
		return
	} else if body.Type == "file" && body.Filename == "" {
		util.Responses.Error(w, http.StatusBadRequest, "field 'filename' must be present")
		return
	}

	// Normal message
	if body.Type == "message" {
		// Save message
		message := database.Message{
			ChatId: chat.ID,
			SenderId: uid,
			Type: 0,
			Message: body.Message,
			Timestamp: time.Now().UnixNano(),
		}
		db.NewRecord(message)
		db.Create(&message)

		// Associate with chat
		db.Model(&chat).Association("Messages").Append(&message)

		util.Responses.Success(w)
		return
	}

	// Generate key
	rawFileKey := make([]byte, 32)
	if _, err := rand.Read(rawFileKey); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to generate upload key")
		return
	}
	fileKey := base64.URLEncoding.EncodeToString(rawFileKey)

	// Hash download key for storage
	h := sha512.New()
	h.Write(rawFileKey)
	hashedFileKey := hex.EncodeToString(h.Sum(nil))

	// Remove file name if image
	if body.Type == "image" {
		body.Filename = ""
	}

	// Create file upload link
	id := uuid.NewV4().String()
	file := database.File{
		Path: "./uploaded/" + id,
		Filename: body.Filename,
		UUID: id,
		Key: hashedFileKey,
		Used: false,
	}
	db.NewRecord(file)
	db.Create(&file)

	// Create message database entry
	message := database.Message{
		ChatId: chat.ID,
		SenderId: uid,
		Type: 1,
		Message: body.Message,
		FileId: file.ID,
		Timestamp: time.Now().UnixNano(),
	}
	if body.Type == "file" {
		message.Message = ""
		message.Type = 2
	}

	// Save to database
	db.NewRecord(message)
	db.Create(&message)

	// Associate with chat
	db.Model(&chat).Association("Messages").Append(&message)

	util.Responses.SuccessWithData(w, map[string]string{"url": viper.GetString("http.domain") + "/api/files/" + file.UUID + "?key=" + fileKey})
}
