package util

import (
	"encoding/base64"
	"fmt"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

var JWT = jwtClass{}

type jwtClass struct{}

// Validate an authentication token given the signed string.
func (j jwtClass) Validate(tokenString string, db *gorm.DB) (*jwt.Token, error) {
	// Retrieve token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing message: %v", token.Header["alg"])
		} else if _, ok := token.Header["kid"]; !ok {
			return nil, fmt.Errorf("unable to find key id in token")
		}

		// Get signing key from database
		var t database.Token
		db.Where("id = ?", token.Header["kid"]).First(&t)
		if t.SigningKey == "" {
			return nil, fmt.Errorf("unable to find signing key for token: %v", token.Header["kid"])
		}

		// Decode signing key
		signingKey, err := base64.StdEncoding.DecodeString(t.SigningKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decode signing key: %v", err)
		}

		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

// Check specified user and user from JWT are the same
func (j jwtClass) CheckUser(token *jwt.Token, user database.User, db *gorm.DB) (bool, error) {
	// Retrieve token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims format")
	}

	// Retrieve user specified in token
	var tokenUser database.User
	db.Where("id = ?", claims["sub"]).First(&tokenUser)
	if tokenUser.ID == 0 {
		return false, fmt.Errorf("no user exists at id: %s", claims["sub"])
	}

	// Ensure users are the same
	return user.ID == tokenUser.ID, nil
}

func (j jwtClass) Unvalidated(tokenString string) (*jwt.Token, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token, nil
}
