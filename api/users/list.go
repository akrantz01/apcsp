package users

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "users", "remote_address": r.RemoteAddr, "path": "/api/users", "method": "GET"})

	// Validate initial request on query parameters
	if len(r.URL.RawQuery) == 0 || r.URL.Query().Get("username") == "" {
		logger.WithFields(logrus.Fields{"query_length": len(r.URL.RawQuery), "username": r.URL.Query().Get("username")}).Trace("Query parameter username is non-existent")
		util.Responses.Error(w, http.StatusBadRequest, "query parameter 'username' must be present")
		return
	}
	logger.WithField("username", r.URL.Query().Get("username")).Trace("Validated initial request")

	// Set searched user in logger
	logger = logger.WithField("username", r.URL.Query().Get("username"))

	// Find all users like given username
	var users []database.User
	db.Where("username LIKE ?", r.URL.Query().Get("username")+"%").Limit(10).Find(&users)
	logger.WithField("query", r.URL.Query().Get("username")+"%").Trace("Retrieved all users with similar name")

	// Convert to array of usernames and names
	var response []map[string]string
	for _, user := range users {
		response = append(response, map[string]string{"name": user.Name, "username": user.Username})
	}
	logger.Trace("Transformed user object to only name and username")

	if len(response) == 0 {
		util.Responses.SuccessWithData(w, []string{})
		return
	}
	util.Responses.SuccessWithData(w, response)
	logger.WithField("responses", len(response)).Debug("Retrieved list of all users with similar username")
}
