package authentication

import (
	"bytes"
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"gopkg.in/hlandau/passlib.v1"
	"html/template"
	"net/http"
)

func ResetPassword(db *gorm.DB, mail chan *gomail.Message, resetNotificationTemplate *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_address": r.RemoteAddr, "path": "/api/auth/reset-password", "method": "POST"})

		// Validate request on method, headers, and body
		if r.Method != http.MethodPost {
			logger.WithField("method", r.Method).Trace("Invalid request method")
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if r.Header.Get("Content-Type") != "application/json" {
			logger.WithField("content_type", r.Header.Get("Content-Type")).Trace("Invalid content type header")
			util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
			return
		} else if r.Body == nil {
			logger.Trace("Body is required")
			util.Responses.Error(w, http.StatusBadRequest, "request body is required")
			return
		}
		logger.Trace("Validated initial request")

		// Validate body
		var body struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.WithError(err).Trace("Failed to decode JSON body")
			util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON body: "+err.Error())
			return
		} else if len(body.Password) != 64 {
			logger.WithField("password", len(body.Password)).Trace("Invalid password length")
			util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
			return
		}
		logger.Trace("Successfully validated body")

		// Get user id from token
		token, _ := util.JWT.Unvalidated(r.Header.Get("Authorization"))
		userId, _ := util.JWT.UserId(token)

		// Get user from database
		var user database.User
		db.Where("id = ?", userId).First(&user)
		if user.ID == 0 {
			logger.WithField("id", userId).Trace("Specified user does not exist")
			util.Responses.Error(w, http.StatusBadRequest, "specified user in token does not exist")
			return
		}
		logger.WithField("username", user.Username).Trace("Retrieved user from database based on token")

		// Add username to logger
		logger = logger.WithField("username", user.Username)

		// Hash password for security
		hash, err := passlib.Hash(body.Password)
		if err != nil {
			logrus.WithError(err).Error("Failed to hash password")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		logger.Trace("Hashed given password")

		// Set password
		user.Password = hash
		db.Save(&user)
		logger.Trace("User password updated")

		// Revoke token
		db.Delete(database.Token{}, "id = ?", token.Header["kid"])
		logger.Trace("Deleted associated signing key")

		// Render template to string
		stringBuffer := bytes.NewBuffer([]byte{})
		if err := resetNotificationTemplate.Execute(stringBuffer, map[string]string{"name": user.Name}); err != nil {
			logger.WithError(err).Error("Unable to render template")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to generate email")
			return
		}

		// Assemble and send email
		m := gomail.NewMessage()
		m.SetHeader("From", viper.GetString("email.sender"))
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Chat App - Your Password was Reset")
		m.SetBody("text/html", stringBuffer.String())
		mail <- m

		util.Responses.Success(w)
		logger.Debug("Reset user's password")
	}
}
