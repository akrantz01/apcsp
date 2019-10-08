package authentication

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/packr/v2"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
	"time"
)

func ForgotPassword(db *gorm.DB, mail chan *gomail.Message, box *packr.Box) func(w http.ResponseWriter, r *http.Request) {
	// Load email templates
	templateString, err := box.FindString("reset-password.tmpl")
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load reset password template from box")
	}
	resetPasswordTemplate, err := template.New("reset-password-email").Parse(templateString)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to load reset password template")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_address": r.RemoteAddr, "path": "/api/auth/forgot-password", "method": "GET"})

		// Validate request on method and query parameters
		if r.Method != http.MethodGet {
			logger.WithField("method", r.Method).Trace("Invalid request method")
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		} else if len(r.URL.RawQuery) == 0 || r.URL.Query().Get("username") == "" {
			logger.WithField("username", r.URL.Query().Get("username")).Trace("Invalid query parameter")
			util.Responses.Error(w, http.StatusBadRequest, "query parameter 'username' is required")
			return
		}
		logger.Trace("Validated request")

		// Check if user exists
		var user database.User
		db.Where("username = ?", r.URL.Query().Get("username")).First(&user)
		if user.ID == 0 {
			logger.WithField("username", r.URL.Query().Get("username")).Trace("User not found in database")
			util.Responses.Error(w, http.StatusBadRequest, "specified user does not exist")
			return
		}

		// Add username to logger
		logger = logger.WithField("username", user.Username)

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

		// Render template to string
		stringBuffer := bytes.NewBuffer([]byte{})
		if err := resetPasswordTemplate.Execute(stringBuffer, map[string]string{"name": user.Name, "domain": viper.GetString("http.domain"), "token": signed}); err != nil {
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
