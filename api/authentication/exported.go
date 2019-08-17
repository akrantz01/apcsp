package authentication

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Exported routes for authentication.
// This includes the login and logout routes.
func Authentication(db *gorm.DB) func(r *http.Request, w http.ResponseWriter) {
	return func(r *http.Request, w http.ResponseWriter) {
		switch r.Method {
		case http.MethodPost:
			login(r, w, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
