package files

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
)

func get(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "files", "remote_address": r.RemoteAddr, "path": "/api/files/{file}", "method": "GET"})

	// Validate initial request on path
	vars := mux.Vars(r)
	if _, ok := vars["file"]; !ok {
		logger.WithField("file", vars["file"]).Trace("Invalid value for path parameter")
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'file' must be present")
		return
	}
	logger.WithField("file", vars["file"]).Trace("Validated initial request")

	// Set file id in logger
	logger = logger.WithField("file", vars["file"])

	// Check file exists
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
	logger.WithField("uid", uid).Trace("Confirmed requesting user in chat")

	// Get stat of file
	stat, err := os.Stat("./uploaded/" + file.UUID)
	if os.IsNotExist(err) {
		logger.WithError(err).Error("File with id does not exist on disk")
		util.Responses.Error(w, http.StatusInternalServerError, "specified file does not exist on disk")
		return
	}
	logger.Trace("Ensured file existed on disk")

	// Open file to read contents
	f, err := os.Open("./uploaded/" + file.UUID)
	if err != nil {
		logger.WithError(err).Error("Failed to read file from disk")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to read file from disk")
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.WithError(err).Error("Failed to close reading file from disk")
		}
	}()
	logger.Trace("Opened file to read from disk")

	// Get content type of file
	header := make([]byte, 512)
	if _, err := f.Read(header); err != nil {
		logger.WithError(err).Error("Unable to read file header")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to read file header")
		return
	}
	logger.Trace("Read file header")

	// Write headers
	w.Header().Set("Content-Type", http.DetectContentType(header))
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
	w.WriteHeader(http.StatusOK)
	logger.Trace("Set headers and status code on response")

	// Reset read head
	if _, err := f.Seek(0, 0); err != nil {
		logger.WithError(err).Error("Failed to reset read head on file")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to reset file read head")
		return
	}
	logger.Trace("Reset read head on file")

	// Copy to client
	if _, err := io.Copy(w, f); err != nil {
		logger.WithError(err).Error("Failed to copy file data to client")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to copy to client")
		return
	}
	logger.Debug("Sent specified file to client")
}
