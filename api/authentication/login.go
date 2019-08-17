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
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
	"time"
)

func login(r *http.Request, w http.ResponseWriter, db *gorm.DB) {
	// Validate initial request on Content-Type header and body
	if r.Header.Get("Content-Type") == "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}

	// Validate JSON body
	var body struct {
		Username string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Responses.Error(w, http.StatusBadRequest, "unable to decode JSON: "+err.Error())
		return
	} else if body.Username == "" || body.Password == "" {
		util.Responses.Error(w, http.StatusBadRequest, "fields 'username' and 'password' are required")
		return
	}

	// Check if user exists
	var user database.User
	db.Where("username = ?", body.Username).First(&user)
	if user.ID == 0 {
		util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Validate password
	newHash, err := passlib.Verify(body.Password, user.Password)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Update password hash if needed
	if newHash != "" {
		user.Password = newHash
		db.Save(&user)
	}

	// Create signing key for JWT
	signingKey := make([]byte, 128)
	if _, err := rand.Read(signingKey); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to generate JWT signing key")
		return
	}

	// Save key to database
	storedToken := database.Token{
		SigningKey: base64.StdEncoding.EncodeToString(signingKey),
		UserId:     user.ID,
	}
	db.NewRecord(storedToken)
	db.Create(&storedToken)

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   fmt.Sprint(user.ID),
	})
	token.Header["kid"] = storedToken.ID

	// Sign token
	signed, err := token.SignedString(signingKey)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to sign JWT")
		return
	}

	util.Responses.SuccessWithData(w, map[string]string{"token": signed})
}
