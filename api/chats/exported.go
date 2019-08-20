package chats

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Methods pertaining to all chats such as listing and creation
func AllChats(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// TODO: add modification and deletion
