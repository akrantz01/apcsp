package messages

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/akrantz01/apcsp/api/websockets"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Methods pertaining to all messages such as listing and creation
func AllMessages(hub *websockets.Hub, db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			list(w, r, db)

		case http.MethodPost:
			create(w, r, hub, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Methods pertaining to a specific message such as description, deletion, and updating
func SpecificMessage(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			read(w, r, db)

		case http.MethodPut:
			update(w, r, db)

		case http.MethodDelete:
			deleteMethod(w, r, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
