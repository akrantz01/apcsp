package authentication

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Revoke authentication tokens
func Logout(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token w/o validation
		token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
		if err != nil {
			util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
			return
		}

		// Delete token row from database
		var storedToken database.Token
		db.Where("id = ?", token.Header["kid"]).First(&storedToken)
		db.Delete(&storedToken)

		util.Responses.Success(w)
	}
}
