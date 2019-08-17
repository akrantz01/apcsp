package authentication

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

func logout(r *http.Request, w http.ResponseWriter, db *gorm.DB) {
	// Validate initial request on headers
	if r.Header.Get("Authorization") == "" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Authorization' is required")
		return
	}

	// Verify JWT from headers
	token, err := util.JWT.Validate(r.Header.Get("Authorization"), db)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "invalid token: "+err.Error())
		return
	}

	// Delete token row from database
	var storedToken database.Token
	db.Where("id = ?", token.Header["kid"]).First(&storedToken)
	db.Delete(&storedToken)

	util.Responses.Success(w)
}
