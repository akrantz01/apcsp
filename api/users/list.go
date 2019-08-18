package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on query parameters
	if len(r.URL.RawQuery) == 0 || r.URL.Query().Get("username") == "" {
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'username' must be present")
		return
	}

	// Find all users like given username
	var users []database.User
	db.Where("username LIKE ?", r.URL.Query().Get("username")+"%").Limit(10).Find(&users)

	// Convert to array of usernames and names
	var response []map[string]string
	for _, user := range users {
		response = append(response, map[string]string{"name": user.Name, "username": user.Username})
	}

	if len(response) == 0 {
		util.Responses.SuccessWithData(w, []string{})
		return
	}
	util.Responses.SuccessWithData(w, response)
}
