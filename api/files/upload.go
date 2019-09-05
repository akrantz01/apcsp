package files

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"os"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on headers, path parameters and form
	vars := mux.Vars(r)
	if _, ok := vars["file"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'file' must be present")
		return
	} else if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "multipart/form-data" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'multipart/form-data'")
		return
	}

	// Check if file exists
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

	// Ensure not used
	if file.Used {
		util.Responses.Error(w, http.StatusBadRequest, "upload link has already been used")
		return
	}

	// Set max file size
	// 12 << 27 = 1476395088 bytes or ~ 1.47 GB
	if err := r.ParseMultipartForm(12 << 27); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to set max file size")
		return
	}

	// Retrieve  file from form
	f, _, err := r.FormFile("file")
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "form ")
		return
	}

	// Create output file
	outFile, err := os.Create(file.Path)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to create upload file")
		return
	}

	// Copy from form to file
	if _, err := io.Copy(outFile, f); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to write uploaded file")
		return
	}

	// Set as used
	file.Used = true
	db.Save(&file)

	util.Responses.Success(w)
}
