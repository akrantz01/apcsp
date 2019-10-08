package users

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gobuffalo/packr/v2"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
)

// Methods pertaining to all users such as creation
func AllUsers(db *gorm.DB, mail chan *gomail.Message, box *packr.Box) func(w http.ResponseWriter, r *http.Request) {
	// Load email templates
	templateString, err := box.FindString("verification.tmpl")
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load verification template from box")
	}
	emailVerificationTemplate, err := template.New("verification-email").Parse(templateString)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load verification template")
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
