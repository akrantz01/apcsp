package authentication

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
	"time"
)

func ForgotPassword(db *gorm.DB, mail chan *gomail.Message, resetPasswordTemplate *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_address": r.RemoteAddr, "path": "/api/auth/forgot-password", "method": "POST"})

		// Validate initial request on Content-Type header and body
		if r.Header.Get("Content-Type") != "application/json" {
			logger.WithField("content_type", r.Header.Get("Content-Type")).Trace("Invalid content type")
			util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
			return
		} else if r.Body == nil {
			logger.Trace("No request body given")
			util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
			return
		}
		logger.Trace("Validated initial request")

		// Validate JSON body
		var body struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.WithError(err).Trace("Invalid json body")
			util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
			return
		} else if body.Username == "" {
			logger.WithField("username", body.Username).Trace("Field username not given")
			util.Responses.Error(w, http.StatusBadRequest, "field 'username' is required")
			return
		}
		logger.Trace("Validated JSON body")

		// Check if user exists
		var user database.User
		db.Where("username = ?", body.Username).First(&user)
		if user.ID == 0 {
			logger.WithField("username", body.Username).Trace("User not found in database")
			util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
			return
		}

		// Add username to logger
		logger = logger.WithField("username", body.Username)

		// Create signing key for JWT
		signingKey := make([]byte, 128)
		if _, err := rand.Read(signingKey); err != nil {
			logger.WithError(err).Error("Unable to generate JWT signing key")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to generate JWT signing key")
			return
		}
		logger.Trace("Generated new signing key bytes for JWT")

		// Save key to database
		storedToken := database.Token{
			SigningKey: base64.StdEncoding.EncodeToString(signingKey),
			UserId:     user.ID,
			Type:       database.TokenResetPassword,
		}
		db.NewRecord(storedToken)
		db.Create(&storedToken)
		logger.Trace("Stored signing key and user id in database")

		// Generate token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + tokenExpiration,
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   fmt.Sprint(user.ID),
		})
		token.Header["kid"] = storedToken.ID
		logger.Trace("Generated JWT claims")

		// Sign token
		signed, err := token.SignedString(signingKey)
		if err != nil {
			logger.WithError(err).Error("Unable to sign JWT")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to sign JWT")
			return
		}
		logger.Trace("Signed JWT with signing key")

		// Assemble values for templates
		templateVars := struct {
			Name   string
			Domain string
			Token  string
		}{
			Name:   user.Name,
			Domain: viper.GetString("http.domain"),
			Token:  signed,
		}

		// Render template to string
		stringBuffer := bytes.NewBuffer([]byte{})
		if err := resetPasswordTemplate.Execute(stringBuffer, templateVars); err != nil {
			logger.WithError(err).Error("Unable to render template")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to generate email")
			return
		}

		// Assemble and send reset email
		m := gomail.NewMessage()
		m.SetHeader("From", viper.GetString("email.sender"))
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Chat App - Reset Your Password")
		m.SetBody("text/html", stringBuffer.String())
		mail <- m

		util.Responses.Success(w)
		logger.Debug("New forgot password request from user")
	}
}
