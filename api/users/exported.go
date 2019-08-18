package users

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Methods pertaining to all users such as creation
func AllUsers(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// Methods pertaining to single users such as reading, updating, and deleting
func SpecificUser(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			read(w, r, db)

		case http.MethodPut:
			update(w, r, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
