package users

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
	"gopkg.in/hlandau/passlib.v1"
	"html/template"
	"net/http"
	"time"
)

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB, mail chan *gomail.Message, emailVerificationTemplate *template.Template) {
	logger := logrus.WithFields(logrus.Fields{"app": "users", "remote_address": r.RemoteAddr, "path": "/api/users", "method": "POST"})

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
		Password string `json:"password"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WithError(err).Trace("Invalid json body")
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Username == "" || body.Password == "" || body.Email == "" || body.Name == "" {
		logger.WithFields(logrus.Fields{"username": body.Username, "password": len(body.Password), "email": body.Email, "name": body.Name}).Trace("Field username, password, email, or name not given")
		util.Responses.Error(w, http.StatusBadRequest, "fields 'username', 'password', 'email', and 'name' are required")
		return
	} else if len(body.Password) != 64 {
		logger.WithField("password", len(body.Password)).Trace("Invalid password length")
		util.Responses.Error(w, http.StatusBadRequest, "field 'password' must be of length 64")
		return
	} else if len(body.Email) < 5 || len(body.Email) > 254 {
		logger.WithField("email", len(body.Email)).Trace("Invalid email length")
		util.Responses.Error(w, http.StatusBadRequest, "field 'email' must be of length between 5 and 254 characters")
		return
	}

	// Check if username is taken
	var user database.User
	db.Where("username = ?", body.Username).First(&user)
	if user.ID != 0 {
		logger.WithFields(logrus.Fields{"username": body.Username, "id": user.ID}).Trace("Username is already taken")
		util.Responses.Error(w, http.StatusConflict, "username is already taken")
		return
	}
	logger.WithField("username", body.Username).Trace("Validated username not taken")

	// Check if email is taken
	db.Where("email = ?", body.Email).First(&user)
	if user.ID != 0 {
		logger.WithFields(logrus.Fields{"email": body.Email, "id": user.ID}).Trace("Email is already taken")
		util.Responses.Error(w, http.StatusConflict, "email is already taken")
		return
	}
	logger.WithField("email", body.Email).Trace("Validated email not taken")

	// Hash the password for security
	hash, err := passlib.Hash(body.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	logger.Trace("Hashed given password")

	// Assemble user info from body
	user.Name = body.Name
	user.Email = body.Email
	user.Username = body.Username
	user.Password = hash
	logger.Trace("Assembled user entry in database")

	// Save to database
	db.NewRecord(user)
	db.Save(&user)
	logger.Trace("Created user entry in database")

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
		Type:       database.TokenVerification,
	}
	db.NewRecord(storedToken)
	db.Create(&storedToken)
	logger.Trace("Stored signing key and user id in database")

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 60*60*24*3,
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
	if err := emailVerificationTemplate.Execute(stringBuffer, map[string]string{"name": user.Name, "domain": viper.GetString("http.domain"), "token": signed}); err != nil {
		logger.WithError(err).Error("Unable to render template")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to generate email")
		return
	}

	// Assemble and send verification email
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.sender"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Chat App - Verify Your Email")
	m.SetBody("text/html", stringBuffer.String())
	mail <- m

	util.Responses.Success(w)
	logger.WithFields(logrus.Fields{"username": body.Username, "id": user.ID}).Debug("Created new user")
}
