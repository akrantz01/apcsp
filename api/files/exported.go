package files

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/akrantz01/apcsp/api/websockets"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Methods pertaining to uploading and downloading files/images
func Files(hub *websockets.Hub, db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(w, r, db)

		case http.MethodPost:
			post(w, r, hub, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
