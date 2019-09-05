package files

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func get(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path
	vars := mux.Vars(r)
	if _, ok := vars["file"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'file' must be present")
		return
	}

	// Check file exists
	var file database.File
	db.Where("uuid = ?", vars["file"]).First(&file)
	if file.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "specified file does not exist")
		return
	}

	// Ensure chat exists
	var chat database.Chat
	db.Preload("Users").Where("id = ?", file.ChatId).First(&chat)
	if chat.ID == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "associated chat was deleted")
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
	for _, user := range chat.Users {
		if uid == user.ID {
			valid = true
			break
		}
	}
	if !valid {
		util.Responses.Error(w, http.StatusForbidden, "user is not part of associated chat")
		return
	}

	// Get stat of file
	stat, err := os.Stat("./uploaded/" + file.UUID)
	if os.IsNotExist(err) {
		util.Responses.Error(w, http.StatusInternalServerError, "specified file does not exist on disk")
		return
	}

	// Open file to read contents
	f, err := os.Open("./uploaded/" + file.UUID)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to read file from disk")
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close uploaded file")
		}
	}()

	// Get content type of file
	header := make([]byte, 512)
	if _, err := f.Read(header); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to read file header")
		return
	}

	// Write headers
	w.Header().Set("Content-Type", http.DetectContentType(header))
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Content-Disposition", "attachment; filename=" + file.Filename)
	w.WriteHeader(http.StatusOK)

	// Reset read head
	if _, err := f.Seek(0, 0); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to reset file read head")
		return
	}

	// Copy to client
	if _, err := io.Copy(w, f); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to copy to client")
		return
	}
}
