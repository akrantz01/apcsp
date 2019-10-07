package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strconv"
)

var JWT = jwtClass{}

type jwtClass struct{}

var jwtLogger = logrus.WithField("app", "jwt")

// Validate an authentication token given the signed string.
func (j jwtClass) Validate(tokenString string, tokenType int, db *gorm.DB) (*jwt.Token, error) {
	// Retrieve token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			jwtLogger.WithField("method", token.Header["alg"]).Trace("Invalid signing method for JWT")
			return nil, fmt.Errorf("unexpected signing message: %v", token.Header["alg"])
		} else if _, ok := token.Header["kid"]; !ok {
			jwtLogger.Trace("No key id in token")
			return nil, fmt.Errorf("unable to find key id in token")
		}
		jwtLogger.Trace("Validated token header/parts")

		// Get signing key from database
		var t database.Token
		db.Where("id = ?", token.Header["kid"]).First(&t)
		if t.SigningKey == "" {
			jwtLogger.WithField("key_id", token.Header["kid"]).Trace("No signing key for token")
			return nil, fmt.Errorf("unable to find signing key for token: %v", token.Header["kid"])
		}
		jwtLogger.WithField("key_id", token.Header["kid"]).Trace("Retrieved token signing key from database")

		// Ensure token is of proper type
		if t.Type != uint(tokenType) {
			return nil, errors.New("invalid token type")
		}

		// Decode signing key
		signingKey, err := base64.StdEncoding.DecodeString(t.SigningKey)
		if err != nil {
			jwtLogger.WithError(err).WithField("key_id", token.Header["kid"]).Trace("Unable to decode signing key from database")
			return nil, fmt.Errorf("unable to decode signing key: %v", err)
		}
		jwtLogger.WithField("key_id", token.Header["kid"]).Trace("Decoded signing key")

		return signingKey, nil
	})
	if err != nil {
		jwtLogger.WithError(err).Trace("Invalid signing key")
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		jwtLogger.Trace("Invalid token")
		return nil, fmt.Errorf("token is invalid")
	}
	jwtLogger.Trace("Valid token for given user")

	return token, nil
}

// Check specified user and user from JWT are the same
func (j jwtClass) CheckUser(token *jwt.Token, user database.User, db *gorm.DB) (bool, error) {
	// Retrieve token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		jwtLogger.Trace("Invalid claims format for token")
		return false, fmt.Errorf("invalid claims format")
	}
	jwtLogger.Trace("Retrieved token claims")

	// Retrieve user specified in token
	var tokenUser database.User
	db.Where("id = ?", claims["sub"]).First(&tokenUser)
	if tokenUser.ID == 0 {
		jwtLogger.WithFields(logrus.Fields{"token_uid": claims["sub"], "specified_uid": user.ID}).Trace("No user found for token id")
		return false, fmt.Errorf("no user exists at id: %s", claims["sub"])
	}
	jwtLogger.WithField("uid", claims["sub"]).Trace("Got user id from token")

	// Ensure users are the same
	return user.ID == tokenUser.ID, nil
}

// Get the parts of a token with out any validation
func (j jwtClass) Unvalidated(tokenString string) (*jwt.Token, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		jwtLogger.WithError(err).Trace("Failed to parse token string")
		return nil, err
	}
	jwtLogger.Trace("Parsed token without verification")

	return token, nil
}

// Get the map claims from the token
func (j jwtClass) Claims(token *jwt.Token) jwt.MapClaims {
	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		jwtLogger.Trace("Invalid claims format in token")
		return nil
	} else {
		jwtLogger.Trace("Got token claims")
		return claims
	}
}

// Get user id from token
func (j jwtClass) UserId(token *jwt.Token) (uint, error) {
	// Ensure id is of correct type
	idStr := JWT.Claims(token)["sub"]
	if _, ok := idStr.(string); !ok {
		jwtLogger.Trace("No subject field in token claims")
		return 0, fmt.Errorf("invalid type for 'subject' in token")
	}
	jwtLogger.WithField("id_string", idStr).Trace("Got user id from token as string")

	// Ensure id is a number
	id, err := strconv.ParseUint(idStr.(string), 10, 32)
	if err != nil {
		jwtLogger.WithError(err).WithField("id", idStr).Trace("Token subject was not an integer")
		return 0, fmt.Errorf("invalid type for 'subject' in token")
	}
	jwtLogger.WithField("id", id).Trace("Parsed subject string in token claims")

	return uint(id), nil
}
