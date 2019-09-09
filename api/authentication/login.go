package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
	"time"
)

const tokenExpiration = 60 * 60 * 24 * 3

// Generate authentication tokens for users given a username and password
func Login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{"app": "authentication", "remote_address": r.RemoteAddr, "path": "/api/auth/login", "method": "POST"})

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
			Username string
			Password string
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.WithError(err).Trace("Invalid json body")
			util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
			return
		} else if body.Username == "" || body.Password == "" {
			logger.WithFields(logrus.Fields{"username": body.Username, "password": len(body.Password)}).Trace("Field username or password not given")
			util.Responses.Error(w, http.StatusBadRequest, "fields 'username' and 'password' are required")
			return
		} else if len(body.Password) != 64 {
			logger.WithField("password", len(body.Password)).Trace("Invalid password length")
			util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
			return
		}
		logger.Trace("Validated JSON body")

		// Check if user exists
		var user database.User
		db.Where("username = ?", body.Username).First(&user)
		if user.ID == 0 {
			logger.WithField("username", body.Username).Trace("Username not found in database")
			util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}
		logger.WithField("username", body.Username).Trace("Got user object from database")

		// Add username to logger
		logger = logger.WithField("username", body.Username)

		// Validate password
		newHash, err := passlib.Verify(body.Password, user.Password)
		if err != nil {
			logger.Trace("Invalid password for user")
			util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}
		logger.Trace("Validated password")

		// Update password hash if needed
		if newHash != "" {
			logger.Trace("New hash generated for user")
			user.Password = newHash
			db.Save(&user)
		}

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

		util.Responses.SuccessWithData(w, map[string]string{"token": signed})
		logger.Debug("New login from user")
	}
}
