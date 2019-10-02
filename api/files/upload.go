package files

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/akrantz01/apcsp/api/websockets"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request, hub *websockets.Hub, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "files", "remote_address": r.RemoteAddr, "path": "/api/files/{file}", "method": "POST"})

	// Validate initial request on headers, path parameters and form
	vars := mux.Vars(r)
	if _, ok := vars["file"]; !ok {
		logger.WithField("file", vars["file"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'file' must be present")
		return
	} else if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "multipart/form-data" {
		logger.WithFields(logrus.Fields{"file": vars["file"], "content_type": r.Header.Get("Content-Type")}).Trace("Invalid content type for request")
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'multipart/form-data'")
		return
	}
	logger.WithField("file", vars["file"]).Trace("Validated initial request")

	// Set file id in logger
	logger = logger.WithField("file", vars["file"])

	// Check if file exists
	var file database.File
	db.Where("uuid = ?", vars["file"]).First(&file)
	if file.ID == 0 {
		logger.Trace("File does not exist")
		util.Responses.Error(w, http.StatusBadRequest, "specified file does not exist")
		return
	}
	logger.Trace("Retrieved file from database")

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("id = ?", file.ChatId).First(&chat)
	if chat.ID == 0 {
		logger.WithField("chat", file.ChatId).Trace("Chat associated with file was deleted")
		util.Responses.Error(w, http.StatusBadRequest, "associated chat was deleted")
		return
	}
	logger.WithField("chat", file.ChatId).Trace("Retrieved chat associated with file")

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
		util.Responses.Error(w, http.StatusForbidden, "user is not part of associated chat")
		return
	}

	// Ensure not used
	if file.Used {
		logger.Trace("File has already been uploaded")
		util.Responses.Error(w, http.StatusBadRequest, "upload link has already been used")
		return
	}

	// Get message associated with file
	var message database.Message
	db.Preload("Sender").Where("file_id = ?", file.ID).First(&message)
	if message.ID == 0 {
		logger.Trace("Message associated with file was deleted")
		util.Responses.Error(w, http.StatusBadRequest, "associated message was deleted")
		return
	}
	logger.WithField("message", message.ID).Trace("Retrieved message associated with file")

	// Set max file size
	// 12 << 27 = 1476395088 bytes or ~ 1.47 GB
	if err := r.ParseMultipartForm(12 << 27); err != nil {
		logger.WithError(err).Error("Unable to set max file size to ~1.5 GB")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to set max file size")
		return
	}
	logger.Trace("Set max file size to ~1.5 GB")

	// Retrieve  file from form
	f, h, err := r.FormFile("file")
	if err != nil {
		logger.WithError(err).Trace("Failed to retrieve file from form")
		util.Responses.Error(w, http.StatusBadRequest, "form field 'file' must be present")
		return
	}
	logger.Trace("Retrieved file from multipart form")

	// Ensure proper file content type for images
	if message.Type == database.MessageImage {
		// Check has content type
		if h.Header.Get("Content-Type") == "" {
			logger.WithField("content_type", h.Header.Get("Content-Type")).Trace("Could not determine file content type")
			util.Responses.Error(w, http.StatusBadRequest, "not content type provided for file")
			return
		}
		logger.Trace("Retrieved content type from file")

		// Ensure image
		if !strings.HasPrefix(h.Header.Get("Content-Type"), "image/") {
			logger.WithField("content_type", h.Header.Get("Content-Type")).Trace("Uploaded file is not an image")
			util.Responses.Error(w, http.StatusBadRequest, "uploaded file must be an image")
			return
		}
		logger.Trace("Ensured file was an image")
	}

	// Create output file
	outFile, err := os.Create(file.Path)
	if err != nil {
		logger.WithError(err).Trace("Failed to create output file to upload to")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to create upload file")
		return
	}
	logger.Trace("Created output file for uploading")

	// Copy from form to file
	if _, err := io.Copy(outFile, f); err != nil {
		logger.WithError(err).Trace("Failed to write uploaded file to output file")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to write uploaded file")
		return
	}
	logger.Trace("Wrote uploaded file to disk")

	// Set as used
	file.Used = true
	db.Save(&file)
	logger.Trace("Set file as already uploaded")

	// Push the message over websockets
	for _, user := range chat.Users {
		// Ignore sending user
		if user.ID == uid {
			continue
		}

		// Send message
		hub.PushMessage(user.Username, message, chat.UUID)
	}

	util.Responses.Success(w)
	logger.Debug("Uploaded specified file for message")
}
