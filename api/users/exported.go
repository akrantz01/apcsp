package users

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
)

// Methods pertaining to all users such as creation
func AllUsers(db *gorm.DB, mail chan *gomail.Message) func(w http.ResponseWriter, r *http.Request) {
	// Load email templates
	emailVerificationTemplate, err := template.ParseFiles("templates/verification.tmpl")
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load reset notification template")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			list(w, r, db)

		case http.MethodPost:
			create(w, r, db, mail, emailVerificationTemplate)

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

		case http.MethodDelete:
			deleteMethod(w, r, db)

		default:
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
