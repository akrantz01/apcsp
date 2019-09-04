package files

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
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
	} else if r.URL.Query().Get("key") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "query parameter 'key' must be present")
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

	// Validate key
	keyBytes, err := base64.URLEncoding.DecodeString(r.URL.Query().Get("key"))
	if err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "failed to decode key")
		return
	}

	// Hash for comparison
	h := sha512.New()
	h.Write(keyBytes)
	hashed := hex.EncodeToString(h.Sum(nil))

	// Compare
	if hashed != file.Key {
		util.Responses.Error(w, http.StatusUnauthorized, "invalid upload key")
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
