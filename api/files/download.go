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
	"log"
	"net/http"
	"os"
	"strconv"
)

func get(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on path and query parameters
	vars := mux.Vars(r)
	if _, ok := vars["file"]; !ok {
		util.Responses.Error(w, http.StatusBadRequest, "path parameter 'file' must be present")
		return
	} else if r.URL.Query().Get("key") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "query parameter 'key' must be present")
		return
	}

	// Check file exists
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
		util.Responses.Error(w, http.StatusUnauthorized, "invalid download key")
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
